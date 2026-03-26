<script lang="ts">
  import { api } from '../lib/api';
  import type { AuthResponse } from '../lib/types';

  let { onregister, onlogin }: {
    onregister?: (data: AuthResponse) => void;
    onlogin?: () => void;
  } = $props();

  let username = $state('');
  let password = $state('');
  let name = $state('');
  let error = $state('');
  let loading = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;

    try {
      const data: AuthResponse = await api.register(username, password, name);
      api.setToken(data.access_token);
      onregister?.(data);
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen flex">
  <!-- Brand panel -->
  <div class="hidden lg:flex lg:w-2/5 bg-gradient-to-br from-green-700 via-green-800 to-green-900 flex-col items-center justify-center p-12 text-white">
    <div class="mb-8">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 opacity-90" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
      </svg>
    </div>
    <h1 class="text-4xl font-bold mb-3 tracking-tight">Buku</h1>
    <p class="text-lg text-white/75 text-center leading-relaxed max-w-xs">
      Your personal space for saving and organizing the web.
    </p>
    <div class="mt-12 flex flex-col gap-3 w-full max-w-xs">
      {#each ['Save links in one place', 'Organize by category', 'Search instantly'] as feature}
        <div class="flex items-center gap-3 text-white/80 text-sm">
          <div class="w-5 h-5 rounded-full bg-white/20 flex items-center justify-center shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          {feature}
        </div>
      {/each}
    </div>
  </div>

  <!-- Form panel -->
  <div class="flex-1 flex items-center justify-center p-8 bg-base-100">
    <div class="w-full max-w-sm">
      <!-- Mobile logo -->
      <div class="flex items-center gap-2 mb-8 lg:hidden">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-green-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
        <span class="text-xl font-bold">Buku</span>
      </div>

      <h2 class="text-2xl font-bold text-base-content mb-1">Create an account</h2>
      <p class="text-base-content/60 text-sm mb-8">Get started — it's free</p>

      {#if error}
        <div class="alert alert-error mb-5 text-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{error}</span>
        </div>
      {/if}

      <form onsubmit={handleSubmit} class="flex flex-col gap-4">
        <label class="form-control">
          <div class="label pb-1">
            <span class="label-text font-medium">Full name</span>
          </div>
          <input
            type="text"
            id="name"
            bind:value={name}
            class="input input-bordered focus:input-primary w-full"
            placeholder="Your full name"
            autocomplete="name"
            required
          />
        </label>

        <label class="form-control">
          <div class="label pb-1">
            <span class="label-text font-medium">Username</span>
          </div>
          <input
            type="text"
            id="username"
            bind:value={username}
            class="input input-bordered focus:input-primary w-full"
            placeholder="Choose a username"
            autocomplete="username"
            required
          />
        </label>

        <label class="form-control">
          <div class="label pb-1">
            <span class="label-text font-medium">Password</span>
          </div>
          <input
            type="password"
            id="password"
            bind:value={password}
            class="input input-bordered focus:input-primary w-full"
            placeholder="Choose a password"
            autocomplete="new-password"
            required
          />
        </label>

        <button type="submit" class="btn btn-primary mt-2 w-full" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Creating account…
          {:else}
            Create account
          {/if}
        </button>
      </form>

      <p class="text-center text-sm text-base-content/60 mt-6">
        Already have an account?
        <button class="link link-primary font-medium" onclick={() => onlogin?.()}>
          Sign in
        </button>
      </p>
    </div>
  </div>
</div>
