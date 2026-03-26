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
  let description = $state('');
  let selectedCategoryId = $state('');
  let loading = $state(false);
  let error = $state('');

  $effect(() => {
    if (editingUrl) {
      url = editingUrl.url;
      description = editingUrl.description || '';
      selectedCategoryId = String(editingUrl.category_id);
    } else {
      url = '';
      description = '';
      selectedCategoryId = categories[0]?.id ? String(categories[0].id) : '';
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!url.trim() || !selectedCategoryId) return;

    loading = true;
    error = '';

    try {
      if (editingUrl) {
        const updated: Bookmark = await api.updateURL(
          editingUrl.id,
          url,
          parseInt(selectedCategoryId),
          description || undefined
        );
        onupdate?.(updated);
      } else {
        const created: Bookmark = await api.createURL(
          url,
          parseInt(selectedCategoryId),
          description || undefined
        );
        oncreate?.(created);
        url = '';
        description = '';
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
      <div class="flex-1">
        <input
          type="url"
          id="url"
          bind:value={url}
          class="input input-bordered w-full focus:input-primary"
          placeholder="https://example.com"
          required
        />
      </div>

      <div class="flex-1">
        <input
          type="text"
          id="description"
          bind:value={description}
          class="input input-bordered w-full focus:input-primary"
          placeholder="Description (optional)"
        />
      </div>

      <div class="lg:w-44">
        <select
          id="category"
          bind:value={selectedCategoryId}
          class="select select-bordered w-full focus:select-primary"
          required
        >
          <option value="" disabled>Category</option>
          {#each categories as category}
            <option value={category.id}>{category.name}</option>
          {/each}
        </select>
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
