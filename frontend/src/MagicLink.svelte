<script>
  import { onMount } from 'svelte';

  export let params = {};

  const API_URL = import.meta.env.VITE_API_URL;

  let status = 'loading';  // loading | ready | submitting | done | error
  let appointment = null;
  let result = null;
  let errorMsg = null;

  onMount(async () => {
    try {
      const res = await fetch(`${API_URL}/api/appointments/${params.id}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      appointment = await res.json();
      status = 'ready';
    } catch (err) {
      errorMsg = 'Could not load appointment details.';
      status = 'error';
    }
  });

  async function respond(newStatus) {
    status = 'submitting';
    try {
      const res = await fetch(`${API_URL}/api/appointments/${params.id}/status`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: newStatus }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      result = newStatus;
      status = 'done';
    } catch (err) {
      errorMsg = err.message;
      status = 'error';
    }
  }
</script>

<main class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
  <div class="bg-white shadow-lg rounded-2xl p-8 max-w-sm w-full text-center">

    {#if status === 'loading'}
      <p class="text-gray-400 py-8">Loading appointment...</p>

    {:else if status === 'ready' || status === 'submitting'}
      <h1 class="text-xl font-bold text-gray-800 mb-1">Your Appointment</h1>
      <p class="text-gray-500 text-sm mb-6">Please confirm or cancel</p>

      <div class="bg-gray-50 rounded-lg p-4 mb-6 text-left space-y-2">
        <p class="text-sm text-gray-600">
          <span class="font-medium text-gray-800">Service:</span> {appointment.service_name}
        </p>
        <p class="text-sm text-gray-600">
          <span class="font-medium text-gray-800">When:</span>
          {new Date(appointment.appointment_time).toLocaleDateString('en-US', {
            weekday: 'long', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit'
          })}
        </p>
      </div>

      <div class="space-y-3">
        <button
          on:click={() => respond('confirmed')}
          disabled={status === 'submitting'}
          class="w-full py-3 rounded-xl text-white font-semibold text-lg
                 bg-green-500 hover:bg-green-600 disabled:bg-green-300 transition-colors"
        >
          Confirm Appointment
        </button>
        <button
          on:click={() => respond('cancelled')}
          disabled={status === 'submitting'}
          class="w-full py-3 rounded-xl text-white font-semibold text-lg
                 bg-red-500 hover:bg-red-600 disabled:bg-red-300 transition-colors"
        >
          Cancel Appointment
        </button>
      </div>

    {:else if status === 'done'}
      {#if result === 'confirmed'}
        <div class="py-6">
          <div class="text-5xl mb-4">&#10003;</div>
          <h2 class="text-xl font-bold text-green-700 mb-2">Appointment Confirmed</h2>
          <p class="text-gray-500 text-sm">Thank you! We look forward to seeing you.</p>
        </div>
      {:else}
        <div class="py-6">
          <div class="text-5xl mb-4">&#10007;</div>
          <h2 class="text-xl font-bold text-red-700 mb-2">Appointment Cancelled</h2>
          <p class="text-gray-500 text-sm">Your appointment has been cancelled. You can rebook anytime.</p>
        </div>
      {/if}

    {:else if status === 'error'}
      <div class="py-6">
        <h2 class="text-lg font-bold text-red-700 mb-2">Something went wrong</h2>
        <p class="text-sm text-red-600">{errorMsg}</p>
      </div>
    {/if}

  </div>
</main>
