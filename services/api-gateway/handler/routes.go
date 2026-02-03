package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/taxxy/api-gateway/config"
	authpb "github.com/taxxy/shared/proto/auth"
	driverpb "github.com/taxxy/shared/proto/driver"
	paymentpb "github.com/taxxy/shared/proto/payment"
	riderpb "github.com/taxxy/shared/proto/rider"
	sharedpb "github.com/taxxy/shared/proto/shared"
	trippb "github.com/taxxy/shared/proto/trip"
)

// ─── Gateway struct ──────────────────────────────────────────
// Holds the config so every handler can dial the right service.

type Gateway struct {
	Cfg config.Config
}

func NewGateway(cfg config.Config) *Gateway {
	return &Gateway{Cfg: cfg}
}

// ════════════════════════════════════════════════════════════
//  AUTH ENDPOINTS
// ════════════════════════════════════════════════════════════

// POST /api/auth/register
func (g *Gateway) Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"     binding:"required"`
		Role     string `json:"role"     binding:"required"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.AuthServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := authpb.NewAuthServiceClient(conn).Register(c.Request.Context(), &authpb.RegisterRequest{
		Email: body.Email, Password: body.Password, Name: body.Name, Role: body.Role, Phone: body.Phone,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId, "token": resp.Token})
}

// POST /api/auth/login
func (g *Gateway) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.AuthServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := authpb.NewAuthServiceClient(conn).Login(c.Request.Context(), &authpb.LoginRequest{
		Email: body.Email, Password: body.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId, "token": resp.Token, "role": resp.Role})
}

// ════════════════════════════════════════════════════════════
//  RIDER ENDPOINTS  (all require auth)
// ════════════════════════════════════════════════════════════

// POST /api/rider/profile
func (g *Gateway) CreateRiderProfile(c *gin.Context) {
	userID := GetUserID(c)
	var body struct {
		Name          string `json:"name"            binding:"required"`
		Phone         string `json:"phone"`
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.RiderServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := riderpb.NewRiderServiceClient(conn).CreateProfile(c.Request.Context(), &riderpb.CreateRiderProfileRequest{
		UserId: userID, Name: body.Name, Phone: body.Phone, PaymentMethod: body.PaymentMethod,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"rider_id": resp.RiderId, "status": resp.Status})
}

// GET /api/rider/profile
func (g *Gateway) GetRiderProfile(c *gin.Context) {
	conn := g.dial(g.Cfg.RiderServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := riderpb.NewRiderServiceClient(conn).GetProfile(c.Request.Context(), &riderpb.GetRiderProfileRequest{
		UserId: GetUserID(c),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/rider/request-ride
func (g *Gateway) RequestRide(c *gin.Context) {
	var body struct {
		Pickup      struct{ Lat, Lng float64 } `json:"pickup"       binding:"required"`
		Dropoff     struct{ Lat, Lng float64 } `json:"dropoff"      binding:"required"`
		VehicleType string                     `json:"vehicle_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.RiderServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := riderpb.NewRiderServiceClient(conn).RequestRide(c.Request.Context(), &riderpb.RequestRideRequest{
		UserId:      GetUserID(c),
		Pickup:      &sharedpb.Location{Latitude: body.Pickup.Lat, Longitude: body.Pickup.Lng},
		Dropoff:     &sharedpb.Location{Latitude: body.Dropoff.Lat, Longitude: body.Dropoff.Lng},
		VehicleType: body.VehicleType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"trip_id":        resp.TripId,
		"status":         resp.Status,
		"estimated_fare": resp.EstimatedFare,
		"eta":            resp.Eta,
	})
}

// GET /api/rider/trips
func (g *Gateway) GetRiderTrips(c *gin.Context) {
	conn := g.dial(g.Cfg.RiderServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := riderpb.NewRiderServiceClient(conn).GetMyTrips(c.Request.Context(), &riderpb.GetRiderTripsRequest{
		UserId: GetUserID(c), Limit: 20,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trips": resp.Trips, "total": resp.Total})
}

// ════════════════════════════════════════════════════════════
//  DRIVER ENDPOINTS  (all require auth)
// ════════════════════════════════════════════════════════════

// POST /api/driver/profile
func (g *Gateway) CreateDriverProfile(c *gin.Context) {
	userID := GetUserID(c)
	var body struct {
		Name         string `json:"name"          binding:"required"`
		Phone        string `json:"phone"`
		VehicleType  string `json:"vehicle_type"  binding:"required"`
		VehiclePlate string `json:"vehicle_plate" binding:"required"`
		VehicleBrand string `json:"vehicle_brand" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.DriverServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := driverpb.NewDriverServiceClient(conn).CreateProfile(c.Request.Context(), &driverpb.CreateDriverProfileRequest{
		UserId: userID, Name: body.Name, Phone: body.Phone,
		VehicleType: body.VehicleType, VehiclePlate: body.VehiclePlate, VehicleBrand: body.VehicleBrand,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"driver_id": resp.DriverId, "status": resp.Status})
}

// GET /api/driver/profile
func (g *Gateway) GetDriverProfile(c *gin.Context) {
	conn := g.dial(g.Cfg.DriverServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := driverpb.NewDriverServiceClient(conn).GetProfile(c.Request.Context(), &driverpb.GetDriverProfileRequest{
		UserId: GetUserID(c),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// PUT /api/driver/availability
func (g *Gateway) UpdateDriverAvailability(c *gin.Context) {
	var body struct {
		Status   string                      `json:"status"   binding:"required"`
		Location *struct{ Lat, Lng float64 } `json:"location"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.DriverServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	req := &driverpb.UpdateAvailabilityRequest{UserId: GetUserID(c), Status: body.Status}
	if body.Location != nil {
		req.Location = &sharedpb.Location{Latitude: body.Location.Lat, Longitude: body.Location.Lng}
	}

	resp, err := driverpb.NewDriverServiceClient(conn).UpdateAvailability(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": resp.Status})
}

// POST /api/driver/accept-ride
func (g *Gateway) AcceptRide(c *gin.Context) {
	var body struct {
		TripID string `json:"trip_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.DriverServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := driverpb.NewDriverServiceClient(conn).AcceptRide(c.Request.Context(), &driverpb.AcceptRideRequest{
		UserId: GetUserID(c), TripId: body.TripID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trip_id": resp.TripId, "status": resp.Status})
}

// POST /api/driver/complete-ride
func (g *Gateway) CompleteRide(c *gin.Context) {
	var body struct {
		TripID string `json:"trip_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.DriverServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := driverpb.NewDriverServiceClient(conn).CompleteRide(c.Request.Context(), &driverpb.CompleteRideRequest{
		UserId: GetUserID(c), TripId: body.TripID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trip_id": resp.TripId, "status": resp.Status, "fare": resp.Fare})
}

// ════════════════════════════════════════════════════════════
//  TRIP ENDPOINTS
// ════════════════════════════════════════════════════════════

// GET /api/trip/:trip_id
func (g *Gateway) GetTrip(c *gin.Context) {
	tripID := c.Param("trip_id")

	conn := g.dial(g.Cfg.TripServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := trippb.NewTripServiceClient(conn).GetTrip(c.Request.Context(), &trippb.GetTripRequest{TripId: tripID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ════════════════════════════════════════════════════════════
//  PAYMENT ENDPOINTS
// ════════════════════════════════════════════════════════════

// POST /api/payment/calculate-fare
func (g *Gateway) CalculateFare(c *gin.Context) {
	var body struct {
		PickupLat   float64 `json:"pickup_lat"   binding:"required"`
		PickupLng   float64 `json:"pickup_lng"   binding:"required"`
		DropoffLat  float64 `json:"dropoff_lat"  binding:"required"`
		DropoffLng  float64 `json:"dropoff_lng"  binding:"required"`
		VehicleType string  `json:"vehicle_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.PaymentServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := paymentpb.NewPaymentServiceClient(conn).CalculateFare(c.Request.Context(), &paymentpb.CalculateFareRequest{
		PickupLat: body.PickupLat, PickupLng: body.PickupLng,
		DropoffLat: body.DropoffLat, DropoffLng: body.DropoffLng,
		VehicleType: body.VehicleType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"base_fare":        resp.BaseFare,
		"distance_km":      resp.DistanceKm,
		"duration_min":     resp.DurationMin,
		"surge_multiplier": resp.SurgeMultiplier,
		"total_fare":       resp.TotalFare,
		"breakdown":        resp.Breakdown,
	})
}

// POST /api/payment/pay
func (g *Gateway) ProcessPayment(c *gin.Context) {
	var body struct {
		TripID        string  `json:"trip_id"         binding:"required"`
		Amount        float64 `json:"amount"          binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := g.dial(g.Cfg.PaymentServiceURL, c)
	if conn == nil {
		return
	}
	defer conn.Close()

	resp, err := paymentpb.NewPaymentServiceClient(conn).ProcessPayment(c.Request.Context(), &paymentpb.ProcessPaymentRequest{
		TripId: body.TripID, RiderId: GetUserID(c), Amount: body.Amount, PaymentMethod: body.PaymentMethod,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_id": resp.PaymentId, "status": resp.Status, "amount": resp.Amount})
}

// ─── dial helper ─────────────────────────────────────────────

func (g *Gateway) dial(url string, c *gin.Context) *grpc.ClientConn {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service unavailable"})
		return nil
	}
	return conn
}
