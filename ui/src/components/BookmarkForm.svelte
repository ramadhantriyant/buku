<script lang="ts">
  import { api } from '../lib/api';
  import type { Bookmark, Category } from '../lib/types';

  let { categories = [], editingUrl = null, oncreate, onupdate, oncancel }: {
    categories?: Category[];
    editingUrl?: Bookmark | null;
    oncreate?: (bookmark: Bookmark) => void;
    onupdate?: (bookmark: Bookmark) => void;
    oncancel?: () => void;
  } = $props();

  let url = $state('');
  let title = $state('');
  let description = $state('');
  let isPinned = $state(false);
  let selectedCategoryId = $state('');
  let loading = $state(false);
  let error = $state('');

  $effect(() => {
    if (editingUrl) {
      url = editingUrl.url;
      title = editingUrl.title || '';
      description = editingUrl.description || '';
      isPinned = editingUrl.is_pinned;
      selectedCategoryId = editingUrl.category_id ? String(editingUrl.category_id) : '';
    } else {
      url = '';
      title = '';
      description = '';
      isPinned = false;
      selectedCategoryId = '';
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!url.trim()) return;

    loading = true;
    error = '';

    const categoryId = selectedCategoryId ? parseInt(selectedCategoryId) : undefined;

    try {
      if (editingUrl) {
        const updated: Bookmark = await api.updateURL(
          editingUrl.id,
          url,
          title || undefined,
          description || undefined,
          isPinned,
          categoryId
        );
        onupdate?.(updated);
      } else {
        const created: Bookmark = await api.createURL(
          url,
          title || undefined,
          description || undefined,
          isPinned,
          categoryId
        );
        oncreate?.(created);
        url = '';
        title = '';
        description = '';
        isPinned = false;
        selectedCategoryId = '';
      }
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="bg-base-100 rounded-xl border border-base-200 shadow-sm mb-6 overflow-hidden">
  <div class="px-5 py-3.5 border-b border-base-200 bg-base-50 flex items-center gap-2">
    {#if editingUrl}
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-green-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
      </svg>
      <span class="text-sm font-semibold text-base-content">Edit bookmark</span>
    {:else}
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-green-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      <span class="text-sm font-semibold text-base-content">Add bookmark</span>
    {/if}
  </div>

  <form onsubmit={handleSubmit} class="p-5">
    {#if error}
      <div class="alert alert-error mb-4 text-sm">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{error}</span>
      </div>
    {/if}

    <div class="flex flex-col lg:flex-row gap-3">
      <!-- URL -->
      <div class="flex-1">
        <input
          type="url"
          bind:value={url}
          class="input input-bordered w-full focus:input-primary"
          placeholder="https://example.com"
          required
        />
      </div>

      <!-- Title -->
      <div class="flex-1">
        <input
          type="text"
          bind:value={title}
          class="input input-bordered w-full focus:input-primary"
          placeholder="Title (optional)"
        />
      </div>

      <!-- Description -->
      <div class="flex-1">
        <input
          type="text"
          bind:value={description}
          class="input input-bordered w-full focus:input-primary"
          placeholder="Description (optional)"
        />
      </div>

      <!-- Category (optional) -->
      <div class="lg:w-44">
        <select
          bind:value={selectedCategoryId}
          class="select select-bordered w-full focus:select-primary"
        >
          <option value="">No category</option>
          {#each categories as category}
            <option value={category.id}>{category.name}</option>
          {/each}
        </select>
      </div>

      <!-- Pin toggle -->
      <div class="flex items-center gap-2 shrink-0">
        <label class="flex items-center gap-1.5 cursor-pointer select-none text-sm text-base-content/70">
          <input type="checkbox" bind:checked={isPinned} class="checkbox checkbox-sm checkbox-primary" />
          Pin
        </label>
      </div>

      <div class="flex gap-2 shrink-0">
        {#if editingUrl}
          <button type="button" class="btn btn-ghost" onclick={() => oncancel?.()}>
            Cancel
          </button>
        {/if}
        <button type="submit" class="btn btn-primary" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
          {:else if editingUrl}
            Update
          {:else}
            Save
          {/if}
        </button>
      </div>
    </div>
  </form>
</div>
