<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import { Input } from "$lib/components/ui/input";
  import { SampleProduct } from "$lib/components/product";
  import { AppStore } from "$lib/components/app_store";
  import {
    Table,
    TableHeader,
    TableHead,
    TableBody,
    TableRow,
    TableCell,
  } from "$lib/components/ui/table";
  import { toast } from "svelte-sonner";

  let updateTimeouts = new Map();
  // Store for parent product IDs, since it needs a string value for Select
  let parent_product_id: { [key: string]: string } = $state({});

  AppStore.subscribe((data) => {
    data.products.forEach((product) => {
      if (product.parentProductID) {
        parent_product_id[product.id] = product.parentProductID;
      } else {
        parent_product_id[product.id] = "";
      }
    });
  });

  async function updateProduct(product: SampleProduct) {
    let payload = {
      name: product.name,
      parent_product_id: parent_product_id[product.id] || null,
    };
    const res = await fetch(`/api/product/${product.id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      toast.error("Failed to update product.");
      throw new Error(await res.text());
    }
    toast.success("Product updated successfully.");
  }

  // Debounced update function
  function debounceUpdate(product: SampleProduct) {
    if (updateTimeouts.has(product.id)) {
      clearTimeout(updateTimeouts.get(product.id));
    }
    const timeout = setTimeout(() => {
      updateProduct(product);
      updateTimeouts.delete(product.id);
    }, 1000);
    updateTimeouts.set(product.id, timeout);
  }

  // Add a new product row by creating it on the backend
  async function addProductRow() {
    await fetch("/api/product", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: "", parent_product_id: null }),
    });
  }

  async function deleteProduct(product: SampleProduct) {
    // If the product has children, we cannot delete it
    if (product.ChildProducts.length > 0) {
      toast.error("Cannot delete product with child products.");
      return;
    }
    await fetch(`/api/product/${product.id}`, { method: "DELETE" });
  }
</script>

<Table class="w-full">
  <TableHeader>
    <TableRow>
      <TableHead>Name</TableHead>
      <TableHead>Parent Product</TableHead>
      <TableHead>Delete</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    {#each $AppStore.products as product, i (product.id || i)}
      <TableRow>
        <TableCell>
          <Input
            class="input input-bordered w-full"
            bind:value={product.name}
            oninput={(e: Event) => {
              const target = e.target as HTMLInputElement;
              product.name = target.value;
              debounceUpdate(product);
            }}
            placeholder="Product Name"
          />
        </TableCell>
        <TableCell>
          <ProductSelect
            bind:bindValue={parent_product_id[product.id]}
            onValueChange={(v) => {
              parent_product_id[product.id] = v;
              parent_product_id = parent_product_id;
              debounceUpdate(product);
            }}
            placeholder="None"
            id={`parent-product-${product.id}`}
            disabled={false}
            required={false}
            filterOutIds={[
              product.id,
              ...product.ChildProducts.map((child) => child.id),
            ]}
            options={$AppStore.products
              .filter(
                (p) =>
                  p.id !== product.id &&
                  !product.ChildProducts.some((child) => child.id === p.id)
              )
              .map((p) => ({ value: p.id, label: p.name }))}
          />
        </TableCell>
        <TableCell>
          <Button
            type="button"
            variant="destructive"
            onclick={() => deleteProduct(product)}
          >
            Delete
          </Button>
        </TableCell>
      </TableRow>
    {/each}
  </TableBody>
</Table>
<div class="mt-4 w-full flex flex-row justify-end">
  <Button type="button" onclick={addProductRow}>Add Product</Button>
</div>
