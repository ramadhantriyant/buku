<script lang="ts">
  import { api } from '../lib/api';
  import type { Category } from '../lib/types';

  let { categories: categoriesProp = [], selectedCategory = null, searchQuery = '', onselect, onsearch }: {
    categories?: Category[];
    selectedCategory?: Category | null;
    searchQuery?: string;
    onselect?: (category: Category | null) => void;
    onsearch?: (query: string) => void;
  } = $props();

  let categories = $state<Category[]>([]);
  let localSearchQuery = $state('');

  $effect(() => { categories = [...categoriesProp]; });
  $effect(() => { localSearchQuery = searchQuery; });

  let newCategoryName = $state('');
  let newCategoryDescription = $state('');
  let showAddCategory = $state(false);
  let loading = $state(false);
  let error = $state('');

  async function addCategory() {
    if (!newCategoryName.trim()) return;

    loading = true;
    error = '';

    try {
      const category: Category = await api.createCategory(
        newCategoryName,
        newCategoryDescription || undefined
      );
      categories = [...categories, category];
      newCategoryName = '';
      newCategoryDescription = '';
      showAddCategory = false;
      onselect?.(category);
    } catch (err: any) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function deleteCategory(category: Category) {
    if (!confirm(`Delete category "${category.name}"?`)) return;

    try {
      await api.deleteCategory(category.id);
      categories = categories.filter(c => c.id !== category.id);
      if (selectedCategory?.id === category.id) {
        onselect?.(null);
      }
    } catch (err: any) {
      alert(err.message);
    }
  }

  function selectCategory(category: Category) {
    const next = selectedCategory?.id === category.id ? null : category;
    onselect?.(next);
  }
</script>

<aside class="w-64 shrink-0 flex flex-col bg-base-100 border-r border-base-200 min-h-screen">
  <!-- Logo -->
  <div class="flex items-center gap-2.5 px-5 py-4 border-b border-base-200">
    <div class="w-8 h-8 rounded-lg bg-green-700 flex items-center justify-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
      </svg>
    </div>
    <span class="font-bold text-base-content text-lg">Buku</span>
  </div>

  <!-- Search -->
  <div class="px-4 pt-4 pb-2">
    <label class="input input-sm input-bordered flex items-center gap-2 w-full">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-base-content/40 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        type="text"
        bind:value={localSearchQuery}
        oninput={() => onsearch?.(localSearchQuery)}
        class="grow text-sm"
        placeholder="Search bookmarks…"
      />
    </label>
  </div>

  <!-- Categories -->
  <nav class="flex-1 px-3 py-2 overflow-y-auto">
    <p class="text-xs font-semibold text-base-content/40 uppercase tracking-wider px-2 mb-1">Categories</p>

    <button
      class="w-full flex items-center gap-2.5 px-2 py-2 rounded-lg text-sm transition-colors
        {!selectedCategory
          ? 'bg-green-800 text-white font-medium'
          : 'text-base-content/70 hover:bg-base-200'}"
      onclick={() => onselect?.(null)}
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
      </svg>
      All Bookmarks
    </button>

    {#each categories as category}
      <div class="group flex items-center gap-1">
        <button
          class="flex-1 flex items-center gap-2.5 px-2 py-2 rounded-lg text-sm transition-colors
            {selectedCategory?.id === category.id
              ? 'bg-green-800 text-white font-medium'
              : 'text-base-content/70 hover:bg-base-200'}"
          onclick={() => selectCategory(category)}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
          <span class="truncate">{category.name}</span>
        </button>
        <button
          class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-base-300 text-base-content/40 hover:text-error transition-all"
          aria-label="Delete category"
          onclick={(e) => { e.stopPropagation(); deleteCategory(category); }}
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    {/each}
  </nav>

  <!-- Add Category -->
  <div class="px-4 pb-4 border-t border-base-200 pt-3">
    {#if showAddCategory}
      <div class="flex flex-col gap-2">
        {#if error}
          <p class="text-xs text-error">{error}</p>
        {/if}
        <input
          type="text"
          bind:value={newCategoryName}
          class="input input-sm input-bordered w-full"
          placeholder="Category name"
        />
        <input
          type="text"
          bind:value={newCategoryDescription}
          class="input input-sm input-bordered w-full"
          placeholder="Description (optional)"
        />
        <div class="flex gap-2">
          <button class="btn btn-sm btn-primary flex-1" onclick={addCategory} disabled={loading}>
            {#if loading}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              Add
            {/if}
          </button>
          <button class="btn btn-sm btn-ghost" onclick={() => showAddCategory = false}>
            Cancel
          </button>
        </div>
      </div>
    {:else}
      <button
        class="w-full flex items-center justify-center gap-2 btn btn-sm btn-ghost border border-dashed border-base-300 text-base-content/50 hover:border-green-400 hover:text-green-700"
        onclick={() => showAddCategory = true}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New category
      </button>
    {/if}
  </div>
</aside>
