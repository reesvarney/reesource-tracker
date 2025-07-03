<script lang="ts">
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Select from "$lib/components/ui/select";
  import { Input } from "$lib/components/ui/input";

  import {
    Table,
    TableHeader,
    TableHead,
    TableBody,
    TableRow,
    TableCell,
  } from "$lib/components/ui/table";

  type Product = {
    ID: string;
    Name: string;
    ParentProductID: string | null;
  };

  let products: Array<Product & { CombinedName: string }> = [];
  let editingProduct: Product | null = null;
  let name = "";
  let parentProductId: string = "";
  let isEditing = false;
  let error = "";

  async function fetchProducts() {
    products = (await (await fetch("/api/products")).json()) ?? [];
  }
  // Inline save for table editing
  async function saveProductInline(
    product: Product & { CombinedName: string },
    newName: string,
    newParentProductId: string
  ) {
    let payload = {
      name: newName,
      parent_product_id: newParentProductId || null,
    };
    let resp = await fetch(`/api/product/${product.ID}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!resp.ok) {
      // Optionally show error UI
      error = (await resp.json()).error || "Failed to save product.";
      await fetchProducts();
      return;
    }
    await fetchProducts();
  }

  function startEdit(product: Product) {
    editingProduct = product;
    name = product.Name;
    parentProductId = product.ParentProductID || "";
    isEditing = true;
    error = "";
  }

  function startCreate() {
    editingProduct = null;
    name = "";
    parentProductId = "";
    isEditing = true;
    error = "";
  }

  async function saveProduct() {
    if (!name) {
      error = "Product name is required.";
      return;
    }
    let id = editingProduct ? editingProduct.ID : "";
    let payload = {
      name,
      parent_product_id: parentProductId || null,
    };
    let resp = await fetch(`/api/product/${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!resp.ok) {
      error = (await resp.json()).error || "Failed to save product.";
      return;
    }
    isEditing = false;
    await fetchProducts();
    console.log(products);
  }

  onMount(fetchProducts);
</script>

{#if isEditing}
  <form
    on:submit|preventDefault={saveProduct}
    class="flex flex-col gap-4 max-w-md"
  >
    <div>
      <label for="product-name" class="block mb-1 font-semibold"
        >Product Name</label
      >
      <Input
        id="product-name"
        class="input input-bordered w-full"
        bind:value={name}
        required
      />
    </div>
    <div>
      <label for="parent-product" class="block mb-1 font-semibold"
        >Parent Product</label
      >
      <Select.Root bind:value={parentProductId} class="w-full">
        <Select.Trigger id="parent-product" class="input input-bordered w-full">
          {#if parentProductId}
            {#if products.find((p) => p.ID === parentProductId)}
              {products.find((p) => p.ID === parentProductId).CombinedName}
            {:else}
              None
            {/if}
          {:else}
            None
          {/if}
        </Select.Trigger>
        <Select.Content class="w-full">
          <Select.Item value="">None</Select.Item>
          {#each products.filter((p) => !editingProduct || p.ID !== editingProduct.ID) as product}
            <Select.Item value={product.ID}>{product.CombinedName}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
    {#if error}
      <div class="text-red-600">{error}</div>
    {/if}
    <div class="flex gap-2">
      <Button type="submit">Save</Button>
      <Button type="button" onclick={() => (isEditing = false)}>Cancel</Button>
    </div>
  </form>
{:else}
  <div class="mb-4">
    <Button type="button" onclick={startCreate}>Add Product</Button>
  </div>
  <Table class="w-full">
    <TableHeader>
      <TableRow>
        <TableHead>Name</TableHead>
        <TableHead>Parent Product</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      {#each products as product, i (product.ID)}
        <TableRow>
          <TableCell>
            <Input
              class="input input-bordered w-full"
              value={product.ItemName}
              onchange={async (e) => {
                const newName = e.target.value;
                if (newName !== product.Name) {
                  await saveProductInline(
                    product,
                    newName,
                    product.ParentProductID || ""
                  );
                }
              }}
            />
          </TableCell>
          <TableCell>
            <Select.Root
              bind:value={product.ParentProductID}
              class="w-full"
              onchange={async (e) => {
                const newParent = e.detail;
                if (newParent !== product.ParentProductID) {
                  await saveProductInline(product, product.Name, newParent);
                }
              }}
            >
              <Select.Trigger class="input input-bordered w-full">
                {#if product.ParentProductID}
                  {#if products.find((p) => p.ID === product.ParentProductID)}
                    {products.find((p) => p.ID === product.ParentProductID)
                      .CombinedName}
                  {:else}
                    None
                  {/if}
                {:else}
                  None
                {/if}
              </Select.Trigger>
              <Select.Content class="w-full">
                <Select.Item value="">None</Select.Item>
                {#each products.filter((p) => p.ID !== product.ID) as parent}
                  <Select.Item value={parent.ID}
                    >{parent.CombinedName}</Select.Item
                  >
                {/each}
              </Select.Content>
            </Select.Root>
          </TableCell>
        </TableRow>
      {/each}
    </TableBody>
  </Table>
{/if}
