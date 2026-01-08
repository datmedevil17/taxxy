import { Card, CardContent } from "@/components/ui/card"
import { Car, MapPin } from "lucide-react"
import { Link } from "react-router-dom"

export default function Welcome() {
  return (
    <div className="min-h-screen bg-black text-white flex flex-col items-center justify-center p-4">
      <div className="max-w-md w-full space-y-8 text-center">
        <div className="space-y-2">
          <h1 className="text-4xl font-bold tracking-tighter">Taxxy</h1>
          <p className="text-zinc-400">Move the way you want.</p>
        </div>

        <div className="grid gap-4">
          <Link to="/rider">
            <Card className="hover:bg-zinc-900 transition-colors cursor-pointer border-zinc-800 bg-black text-white group">
              <CardContent className="flex items-center justify-between p-6">
                <div className="flex flex-col items-start gap-1">
                  <span className="font-semibold text-lg">I need a ride</span>
                  <span className="text-sm text-zinc-500 group-hover:text-zinc-400">Get where you're going</span>
                </div>
                <MapPin className="h-8 w-8 text-zinc-500 group-hover:text-white transition-colors" />
              </CardContent>
            </Card>
          </Link>

          <Link to="/driver">
            <Card className="hover:bg-zinc-900 transition-colors cursor-pointer border-zinc-800 bg-black text-white group">
              <CardContent className="flex items-center justify-between p-6">
                <div className="flex flex-col items-start gap-1">
                  <span className="font-semibold text-lg">I want to drive</span>
                  <span className="text-sm text-zinc-500 group-hover:text-zinc-400">Earn on your schedule</span>
                </div>
                <Car className="h-8 w-8 text-zinc-500 group-hover:text-white transition-colors" />
              </CardContent>
            </Card>
          </Link>
        </div>
      </div>
    </div>
  )
}
