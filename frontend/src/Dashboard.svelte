<script>
  import { onMount } from 'svelte';

  const API_URL = import.meta.env.VITE_API_URL;

  // Auth state
  let token = localStorage.getItem('admin_token') || '';
  let loggedIn = !!token;
  let passwordInput = '';
  let authError = '';

  // Appointments state
  let appointments = [];
  let loading = false;
  let fetchError = null;

  // Modal state
  let showModal = false;
  let form = { first_name: '', last_name: '', phone: '', service_name: '', appointment_time: '' };
  let submitting = false;
  let submitError = '';

  function authHeaders() {
    return { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` };
  }

  async function fetchAppointments() {
    loading = true;
    fetchError = null;
    try {
      const res = await fetch(`${API_URL}/api/appointments`, { headers: authHeaders() });
      if (res.status === 401) { logout(); return; }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      appointments = await res.json();
    } catch (err) {
      fetchError = err.message;
    } finally {
      loading = false;
    }
  }

  function login() {
    authError = '';
    if (!passwordInput.trim()) { authError = 'Password is required'; return; }
    token = passwordInput.trim();
    localStorage.setItem('admin_token', token);
    loggedIn = true;
    passwordInput = '';
    fetchAppointments();
  }

  function logout() {
    token = '';
    loggedIn = false;
    localStorage.removeItem('admin_token');
    appointments = [];
  }

  async function createAppointment() {
    submitting = true;
    submitError = '';
    try {
      const res = await fetch(`${API_URL}/api/appointments`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({
          ...form,
          appointment_time: new Date(form.appointment_time).toISOString(),
        }),
      });
      if (res.status === 401) { logout(); return; }
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || `HTTP ${res.status}`);
      }
      showModal = false;
      form = { first_name: '', last_name: '', phone: '', service_name: '', appointment_time: '' };
      await fetchAppointments();
    } catch (err) {
      submitError = err.message;
    } finally {
      submitting = false;
    }
  }

  onMount(() => { if (loggedIn) fetchAppointments(); });

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

{#if !loggedIn}
  <!-- Login Screen -->
  <main class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
    <div class="bg-white shadow-lg rounded-xl p-8 max-w-sm w-full">
      <h1 class="text-xl font-bold text-gray-800 mb-1 text-center">project_abc</h1>
      <p class="text-gray-400 text-sm text-center mb-6">Admin Login</p>

      <form on:submit|preventDefault={login}>
        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Admin Password</label>
        <input
          id="password"
          type="password"
          bind:value={passwordInput}
          class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          placeholder="Enter admin password"
        />
        {#if authError}
          <p class="text-red-600 text-xs mt-1">{authError}</p>
        {/if}
        <button
          type="submit"
          class="mt-4 w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 rounded-lg transition-colors"
        >
          Sign In
        </button>
      </form>
    </div>
  </main>

{:else}
  <!-- Dashboard -->
  <main class="min-h-screen bg-gray-50">
    <header class="bg-white border-b border-gray-200">
      <div class="max-w-4xl mx-auto px-6 py-4 flex items-center justify-between">
        <h1 class="text-xl font-bold text-gray-800">project_abc</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-400">Booking Dashboard</span>
          <button on:click={logout} class="text-sm text-red-500 hover:text-red-700 font-medium">
            Logout
          </button>
        </div>
      </div>
    </header>

    <div class="max-w-4xl mx-auto px-6 py-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-lg font-semibold text-gray-700">Appointments</h2>
        <div class="flex gap-3">
          <button
            on:click={() => (showModal = true)}
            class="text-sm bg-indigo-600 hover:bg-indigo-700 text-white font-medium px-4 py-2 rounded-lg transition-colors"
          >
            + New Appointment
          </button>
          <button
            on:click={fetchAppointments}
            class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
          >
            Refresh
          </button>
        </div>
      </div>

      {#if loading}
        <p class="text-gray-400 text-center py-12">Loading appointments...</p>
      {:else if fetchError}
        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
          <p class="text-sm font-semibold text-red-800">Failed to load appointments</p>
          <p class="text-sm text-red-600 mt-1">{fetchError}</p>
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

  <!-- New Appointment Modal -->
  {#if showModal}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">New Appointment</h3>
          <button on:click={() => (showModal = false)} class="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
        </div>

        <form on:submit|preventDefault={createAppointment} class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="first_name" class="block text-sm font-medium text-gray-700 mb-1">First Name</label>
              <input id="first_name" type="text" bind:value={form.first_name} required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
            </div>
            <div>
              <label for="last_name" class="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
              <input id="last_name" type="text" bind:value={form.last_name} required
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
            </div>
          </div>

          <div>
            <label for="phone" class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
            <input id="phone" type="tel" bind:value={form.phone} required placeholder="+48..."
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          </div>

          <div>
            <label for="service" class="block text-sm font-medium text-gray-700 mb-1">Service Name</label>
            <input id="service" type="text" bind:value={form.service_name} required placeholder="e.g. Haircut & Styling"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          </div>

          <div>
            <label for="datetime" class="block text-sm font-medium text-gray-700 mb-1">Date & Time</label>
            <input id="datetime" type="datetime-local" bind:value={form.appointment_time} required
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" />
          </div>

          {#if submitError}
            <p class="text-red-600 text-sm">{submitError}</p>
          {/if}

          <div class="flex gap-3 pt-2">
            <button type="button" on:click={() => (showModal = false)}
              class="flex-1 border border-gray-300 text-gray-700 font-medium py-2 rounded-lg hover:bg-gray-50 transition-colors">
              Cancel
            </button>
            <button type="submit" disabled={submitting}
              class="flex-1 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white font-medium py-2 rounded-lg transition-colors">
              {submitting ? 'Creating...' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}
{/if}
