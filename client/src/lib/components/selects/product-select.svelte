<script lang="ts">
  import * as Select from "../ui/select";
  import { AppStore } from "$lib/components/app_store";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select a product";
  export let id: string = "product-select";
  export let required: boolean = false;
  export let options: { value: string; label: string }[] = [];
  export let onValueChange: (value: string) => void = () => {};
  export let filterMode: boolean = false;
  export let filterOutIds: string[] = [];
</script>

<Select.Root
  {onValueChange}
  bind:value={bindValue}
  type="single"
  {disabled}
  {required}
>
  <Select.Trigger {id} class="w-full">
    {#if bindValue}
      {#if options.length > 0}
        {options.find((o) => o.value === bindValue)?.label}
      {:else if bindValue.startsWith("any-")}
        {(() => {
          const parentId = bindValue.slice(4);
          const parent = $AppStore.products.find((p) => p.id === parentId);
          return `Any variant of ${parent ? parent.name : ""}`;
        })()}
      {:else}
        {$AppStore.products.find((p) => p.id === bindValue)?.name}
      {/if}
    {:else}
      {placeholder}
    {/if}
  </Select.Trigger>
  <Select.Content class="w-full">
    <Select.Item value="">None</Select.Item>
    {#if options.length > 0}
      {#each options.filter((option) => !filterOutIds.includes(option.value)) as option}
        <Select.Item value={option.value}>{option.label}</Select.Item>
      {/each}
    {:else if filterMode}
      {#each $AppStore.products.filter((p) => !p.parentProductID && !filterOutIds.includes(p.id)) as parent}
        <Select.Label>{parent.name}</Select.Label>
        {#if !filterOutIds.includes("any-" + parent.id)}
          <Select.Item value={"any-" + parent.id}
            >Any variant of {parent.name}</Select.Item
          >
        {/if}
        {#each $AppStore.products.filter((p) => p.parentProductID === parent.id && !filterOutIds.includes(p.id)) as child}
          <Select.Item value={child.id}>{child.name}</Select.Item>
        {/each}
      {/each}
    {:else}
      {#each $AppStore.products.filter((p) => p.parentProductID && !filterOutIds.includes(p.id)) as child}
        <Select.Item value={child.id}>{child.name}</Select.Item>
      {/each}
    {/if}
  </Select.Content>
</Select.Root>
