<script>
  const API_URL = import.meta.env.VITE_API_URL;

  let healthData = null;
  let error = null;
  let loading = false;

  async function checkHealth() {
    loading = true;
    error = null;
    healthData = null;

    try {
      const res = await fetch(`${API_URL}/health`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      healthData = await res.json();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<main class="min-h-screen bg-gray-50 flex items-center justify-center">
  <div class="bg-white shadow-lg rounded-xl p-8 max-w-md w-full text-center">
    <h1 class="text-2xl font-bold text-gray-800 mb-2">project_abc</h1>
    <p class="text-gray-500 mb-6">Booking & SMS notifications for the beauty industry</p>

    <button
      on:click={checkHealth}
      disabled={loading}
      class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white font-medium px-6 py-2 rounded-lg transition-colors"
    >
      {loading ? 'Checking...' : 'Check Backend Health'}
    </button>

    {#if healthData}
      <div class="mt-6 bg-green-50 border border-green-200 rounded-lg p-4 text-left">
        <p class="text-sm font-semibold text-green-800 mb-2">Backend is reachable</p>
        <div class="text-sm text-green-700 space-y-1">
          <p><span class="font-medium">Status:</span> {healthData.status}</p>
          <p><span class="font-medium">Project:</span> {healthData.project}</p>
        </div>
      </div>
    {/if}

    {#if error}
      <div class="mt-6 bg-red-50 border border-red-200 rounded-lg p-4 text-left">
        <p class="text-sm font-semibold text-red-800 mb-1">Connection failed</p>
        <p class="text-sm text-red-600">{error}</p>
      </div>
    {/if}
  </div>
</main>
