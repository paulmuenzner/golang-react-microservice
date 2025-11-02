// app/frontend/public/app/page.tsx
import { api } from '@/lib/api'
import Image from "next/image"

// ✅ Server Component - Fetches on server (no useState needed!)
export default async function Home() {
  console.log('API_URL:', process.env.API_URL)
  console.log('NEXT_PUBLIC_API_URL:', process.env.NEXT_PUBLIC_API_URL)

  // Fetch data on server
  let serviceAHealth = null
  let serviceBHealth = null
  let error = null

  try {
    // Parallel requests for better performance
    const [healthA, healthB] = await Promise.all([
      api.serviceA.getHealth(),
      api.serviceB.getHealth(),
    ])
    
    serviceAHealth = healthA
    serviceBHealth = healthB
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to fetch services'
    console.error('Service health check failed:', err)
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <Image
          className="dark:invert"
          src="/next.svg"
          alt="Next.js logo"
          width={100}
          height={20}
          priority
        />
        
        <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
          <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50">
            Microservices Dashboard
          </h1>
          
          {/* Service Status Display */}
          <div className="w-full max-w-md space-y-4">
            {error ? (
              <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg">
                <p className="text-red-600 dark:text-red-400">
                  ❌ Error: {error}
                </p>
              </div>
            ) : (
              <>
                {/* Service A Status */}
                <div className="p-4 bg-zinc-100 dark:bg-zinc-800 rounded-lg">
                  <h3 className="font-semibold text-zinc-900 dark:text-zinc-50">
                    Service A
                  </h3>
                  <p className={`text-sm ${
                    serviceAHealth?.status === 'OK' 
                      ? 'text-green-600 dark:text-green-400' 
                      : 'text-red-600 dark:text-red-400'
                  }`}>
                    {serviceAHealth?.status === 'OK' ? '✅ Healthy' : '❌ Unhealthy'}
                  </p>
                </div>

                {/* Service B Status */}
                <div className="p-4 bg-zinc-100 dark:bg-zinc-800 rounded-lg">
                  <h3 className="font-semibold text-zinc-900 dark:text-zinc-50">
                    Service B
                  </h3>
                  <p className={`text-sm ${
                    serviceBHealth?.status === 'OK' 
                      ? 'text-green-600 dark:text-green-400' 
                      : 'text-red-600 dark:text-red-400'
                  }`}>
                    {serviceBHealth?.status === 'OK' ? '✅ Healthy' : '❌ Unhealthy'}
                  </p>
                </div>
              </>
            )}
          </div>
          
          <p className="max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400">
            Real-time service health status from your microservices backend.
          </p>
        </div>
        
        <div className="flex flex-col gap-4 text-base font-medium sm:flex-row">
          
           <a className="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-foreground px-5 text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc] md:w-[158px]"
            href="http://localhost:3000"
            target="_blank"
            rel="noopener noreferrer"
          >
            📊 Grafana
          </a>
          
          <a  className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-5 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-[158px]"
            href="http://localhost:8080/api/health"
            target="_blank"
            rel="noopener noreferrer"
          >
            🔧 API Gateway
          </a>
        </div>
      </main>
    </div>
  )
}