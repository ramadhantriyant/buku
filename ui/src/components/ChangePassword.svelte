<script lang="ts">
  import { api } from '../lib/api';

  let { onclose }: { onclose?: () => void } = $props();

  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');
  let loading = $state(false);
  let error = $state('');
  let success = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';

    if (newPassword !== confirmPassword) {
      error = 'New passwords do not match.';
      return;
    }
    if (newPassword.length < 6) {
      error = 'New password must be at least 6 characters.';
      return;
    }

    loading = true;
    try {
      await api.changePassword(currentPassword, newPassword);
      success = true;
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) onclose?.();
  }
</script>

<!-- Backdrop -->
<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
  role="presentation"
  onclick={handleBackdropClick}
>
  <div class="bg-base-100 rounded-2xl shadow-xl w-full max-w-sm mx-4 overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 border-b border-base-200">
      <h2 class="font-semibold text-base-content">Change Password</h2>
      <button
        class="btn btn-ghost btn-sm btn-circle text-base-content/50 hover:text-base-content"
        onclick={onclose}
        aria-label="Close"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Body -->
    <div class="px-6 py-5">
      {#if success}
        <div class="flex flex-col items-center gap-3 py-4 text-center">
          <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <p class="font-medium text-base-content">Password updated successfully.</p>
          <button class="btn btn-sm bg-green-800 hover:bg-green-900 text-white border-0 mt-1" onclick={onclose}>
            Done
          </button>
        </div>
      {:else}
        <form onsubmit={handleSubmit} class="flex flex-col gap-3">
          {#if error}
            <p class="text-sm text-error">{error}</p>
          {/if}

          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-base-content/60">Current password</span>
            <input
              type="password"
              bind:value={currentPassword}
              class="input input-bordered input-sm w-full"
              required
              autocomplete="current-password"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-base-content/60">New password</span>
            <input
              type="password"
              bind:value={newPassword}
              class="input input-bordered input-sm w-full"
              required
              autocomplete="new-password"
            />
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-base-content/60">Confirm new password</span>
            <input
              type="password"
              bind:value={confirmPassword}
              class="input input-bordered input-sm w-full"
              required
              autocomplete="new-password"
            />
          </label>

          <div class="flex gap-2 pt-1">
            <button
              type="submit"
              class="btn btn-sm flex-1 bg-green-800 hover:bg-green-900 text-white border-0"
              disabled={loading}
            >
              {#if loading}
                <span class="loading loading-spinner loading-xs"></span>
              {:else}
                Update password
              {/if}
            </button>
            <button type="button" class="btn btn-sm btn-ghost" onclick={onclose}>
              Cancel
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>
</div>
