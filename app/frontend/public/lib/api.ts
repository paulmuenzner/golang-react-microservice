// app/frontend/public/lib/api.ts

/**
 * API Client with context-aware URL resolution
 * 
 * This client automatically selects the correct API URL based on the execution context:
 * - Browser (Client-side): Uses NEXT_PUBLIC_API_URL (accessible from host machine)
 * - Server (Server-side): Uses API_URL (Docker internal network)
 * 
 * Environment Variables:
 * - NEXT_PUBLIC_API_URL: For browser requests (e.g., http://localhost:8080/api)
 * - API_URL: For server-side requests (e.g., http://gateway:8080/api)
 */

/**
 * Determines the correct API base URL based on execution context
 * @returns API base URL (without trailing slash)
 */
const getApiUrl = (): string => {
  // ✅ DEBUG
  console.log('🔍 DEBUG ENV:', {
    'typeof window': typeof window,
    'API_URL': process.env.API_URL,
    'NEXT_PUBLIC_API_URL': process.env.NEXT_PUBLIC_API_URL,
    'NODE_ENV': process.env.NODE_ENV,
  })
  // Server-side rendering (Next.js server)
  // In Docker: uses internal Docker network (gateway:8080)
  if (typeof window === 'undefined') {
    return process.env.API_URL || 'http://gateway:8080/api'
  }
  
  // Client-side rendering (Browser)
  // Always uses host-accessible URL (localhost:8080)
  return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'
}

const API_BASE = getApiUrl()

/**
 * Generic fetch wrapper with error handling
 */
async function apiFetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE}${endpoint}`
  
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(
        `API Error: ${response.status} ${response.statusText} - ${errorText}`
      )
    }

    return await response.json()
  } catch (error) {
    console.error('API Fetch Error:', error)
    throw error
  }
}

/**
 * Type-safe API client
 */
export const api = {
  /**
   * Gateway health check
   */
  getHealth: () => 
    apiFetch<{ status: string }>('/health'),
  
  /**
   * Service A endpoints
   */
  serviceA: {
    getHealth: () => 
      apiFetch<{ status: string }>('/service-a/health'),
    
    getRoot: () => 
      apiFetch<{ message: string }>('/service-a/'),
  },
  
  /**
   * Service B endpoints
   */
  serviceB: {
    getHealth: () => 
      apiFetch<{ status: string }>('/service-b/health'),
    
    getRoot: () => 
      apiFetch<{ message: string }>('/service-b/'),
  },
}

/**
 * Export base URL for debugging
 */
export const getApiBaseUrl = () => API_BASE