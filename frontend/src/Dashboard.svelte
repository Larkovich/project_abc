<script>
  import { onMount } from 'svelte';

  const API_URL = import.meta.env.VITE_API_URL;

  let appointments = [];
  let loading = true;
  let error = null;

  async function fetchAppointments() {
    loading = true;
    error = null;
    try {
      const res = await fetch(`${API_URL}/api/appointments`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      appointments = await res.json();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  onMount(fetchAppointments);

  function formatDate(iso) {
    return new Date(iso).toLocaleDateString('en-US', {
      weekday: 'short', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
  }

  const statusColors = {
    scheduled: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
    completed: 'bg-gray-100 text-gray-600',
  };
</script>

<main class="min-h-screen bg-gray-50">
  <header class="bg-white border-b border-gray-200">
    <div class="max-w-4xl mx-auto px-6 py-4 flex items-center justify-between">
      <h1 class="text-xl font-bold text-gray-800">project_abc</h1>
      <span class="text-sm text-gray-400">Booking Dashboard</span>
    </div>
  </header>

  <div class="max-w-4xl mx-auto px-6 py-8">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-lg font-semibold text-gray-700">Appointments</h2>
      <button
        on:click={fetchAppointments}
        class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
      >
        Refresh
      </button>
    </div>

    {#if loading}
      <p class="text-gray-400 text-center py-12">Loading appointments...</p>
    {:else if error}
      <div class="bg-red-50 border border-red-200 rounded-lg p-4">
        <p class="text-sm font-semibold text-red-800">Failed to load appointments</p>
        <p class="text-sm text-red-600 mt-1">{error}</p>
      </div>
    {:else if appointments.length === 0}
      <p class="text-gray-400 text-center py-12">No appointments found.</p>
    {:else}
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left text-xs font-medium text-gray-500 uppercase px-6 py-3">Time</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase px-6 py-3">Customer</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase px-6 py-3">Service</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase px-6 py-3">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            {#each appointments as appt}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 text-sm text-gray-700">{formatDate(appt.appointment_time)}</td>
                <td class="px-6 py-4 text-sm text-gray-800 font-medium">
                  {appt.customer_first_name} {appt.customer_last_name}
                </td>
                <td class="px-6 py-4 text-sm text-gray-600">{appt.service_name}</td>
                <td class="px-6 py-4">
                  <span class="text-xs font-medium px-2 py-1 rounded-full {statusColors[appt.status] || 'bg-gray-100 text-gray-600'}">
                    {appt.status}
                  </span>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</main>
