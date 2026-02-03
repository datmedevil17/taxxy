import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Car, MapPin, LogIn } from "lucide-react";
import { Link } from "react-router-dom";
import { authService } from "@/services/auth";

export default function Welcome() {
  const isAuthenticated = authService.isAuthenticated();
  const userRole = authService.getRole();

  return (
    <div className="min-h-screen bg-gradient-to-br from-zinc-900 via-black to-zinc-900 text-white flex flex-col items-center justify-center p-4">
      <div className="max-w-md w-full space-y-8 text-center">
        <div className="space-y-4">
          <img src="/favicon.png" alt="Taxxy" className="h-20 w-auto mx-auto" />
          <h1 className="text-5xl font-bold tracking-tighter bg-gradient-to-r from-white to-zinc-400 bg-clip-text text-transparent">
            Taxxy
          </h1>
          <p className="text-zinc-400">Move the way you want.</p>
        </div>

        {/* Auth Buttons */}
        {!isAuthenticated ? (
          <div className="flex gap-3 justify-center">
            <Link to="/login">
              <Button variant="outline" className="border-zinc-700 text-white hover:bg-zinc-900 gap-2">
                <LogIn className="h-4 w-4" /> Sign In
              </Button>
            </Link>
            <Link to="/register">
              <Button className="bg-white text-black hover:bg-zinc-200">
                Get Started
              </Button>
            </Link>
          </div>
        ) : (
          <div className="text-sm text-zinc-500">
            Welcome back! Choose your mode below.
          </div>
        )}

        <div className="grid gap-4 pt-4">
          <Link to={isAuthenticated && userRole === 'rider' ? '/rider' : '/register'}>
            <Card className="hover:bg-white/10 transition-colors cursor-pointer border-white/10 bg-white/5 backdrop-blur-lg text-white group">
              <CardContent className="flex items-center justify-between p-6">
                <div className="flex flex-col items-start gap-1">
                  <span className="font-semibold text-lg">I need a ride</span>
                  <span className="text-sm text-zinc-400 group-hover:text-zinc-300">Get where you're going</span>
                </div>
                <MapPin className="h-8 w-8 text-zinc-400 group-hover:text-white transition-colors" />
              </CardContent>
            </Card>
          </Link>

          <Link to={isAuthenticated && userRole === 'driver' ? '/driver' : '/register'}>
            <Card className="hover:bg-white/10 transition-colors cursor-pointer border-white/10 bg-white/5 backdrop-blur-lg text-white group">
              <CardContent className="flex items-center justify-between p-6">
                <div className="flex flex-col items-start gap-1">
                  <span className="font-semibold text-lg">I want to drive</span>
                  <span className="text-sm text-zinc-400 group-hover:text-zinc-300">Earn on your schedule</span>
                </div>
                <Car className="h-8 w-8 text-zinc-400 group-hover:text-white transition-colors" />
              </CardContent>
            </Card>
          </Link>
        </div>

        {isAuthenticated && (
          <Button
            variant="ghost"
            className="text-zinc-500 hover:text-white"
            onClick={() => {
              authService.logout();
              window.location.reload();
            }}
          >
            Sign Out
          </Button>
        )}
      </div>
    </div>
  );
}
