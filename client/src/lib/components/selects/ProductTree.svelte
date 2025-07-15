<script lang="ts">
  import * as Select from "$lib/components/ui/select";
  import type { SampleProduct } from "$lib/components/product";
  import ProductTree from "./ProductTree.svelte";
  export let products: SampleProduct[] = [];
  export let filterOutIds: string[] = [];
  export let level: number = 0;
  export let filterMode: boolean = false;
</script>

{#each products as product (product.id)}
  {#if !filterOutIds.includes(product.id)}
    {#if product.ChildProducts && product.ChildProducts.length > 0}
      <Select.Label style="padding-left: max({level * 1.5}em, 0.5em)"
        >{product.CombinedName}</Select.Label
      >
      {#if !filterOutIds.includes(`any-${product.id}`) && product.ChildProducts.length > 1 && filterMode}
        <Select.Item
          value={"any-" + product.id}
          style="padding-left: {(level + 1) * 1.5}em"
        >
          Any variant of {product.CombinedName}
        </Select.Item>
      {/if}
      <ProductTree
        products={product.ChildProducts}
        {filterOutIds}
        level={level + 1}
        {filterMode}
      />
    {:else}
      <Select.Item value={product.id} style="padding-left: {level * 1.5}em"
        >{product.CombinedName}</Select.Item
      >
    {/if}
  {/if}
{/each}
