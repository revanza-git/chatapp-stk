export default function DebugPage() {
  // Get all environment variables that start with NEXT_PUBLIC_
  const envVars = Object.entries(process.env).filter(([key]) => 
    key.startsWith('NEXT_PUBLIC_')
  );

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Environment Variables Debug</h1>
      
      <div className="bg-gray-100 p-4 rounded-lg mb-6">
        <h2 className="text-lg font-semibold mb-4">NEXT_PUBLIC_ Environment Variables:</h2>
        {envVars.length === 0 ? (
          <p className="text-red-600">❌ No NEXT_PUBLIC_ environment variables found!</p>
        ) : (
          <ul className="space-y-2">
            {envVars.map(([key, value]) => (
              <li key={key} className="flex flex-col">
                <span className="font-mono text-sm font-bold">{key}:</span>
                <span className="font-mono text-sm bg-white p-2 rounded border">
                  {value || '(empty)'}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="bg-blue-100 p-4 rounded-lg mb-6">
        <h2 className="text-lg font-semibold mb-4">API URL Being Used:</h2>
        <p className="font-mono text-sm bg-white p-2 rounded border">
          {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080 (fallback)'}
        </p>
      </div>

      <div className="bg-yellow-100 p-4 rounded-lg">
        <h2 className="text-lg font-semibold mb-4">Build Information:</h2>
        <ul className="space-y-1 text-sm">
          <li><strong>NODE_ENV:</strong> {process.env.NODE_ENV}</li>
          <li><strong>Build Time:</strong> {new Date().toISOString()}</li>
        </ul>
      </div>
    </div>
  );
} 