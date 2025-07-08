<script lang="ts">
  import * as Select from "../ui/select";
  import { AppStore } from "$lib/components/app_store";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select a location";
  export let id: string = "location-select";
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
          const parent = $AppStore.locations.find((l) => l.id === parentId);
          return `Anywhere in ${parent ? parent.name : ""}`;
        })()}
      {:else}
        {$AppStore.locations.find((l) => l.id === bindValue)?.name}
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
      {#each $AppStore.locations.filter((l) => !l.parentLocationID && !filterOutIds.includes(l.id)) as parent}
        <Select.Label>{parent.name}</Select.Label>
        {#if !filterOutIds.includes("any-" + parent.id)}
          <Select.Item value={"any-" + parent.id}
            >Anywhere in {parent.name}</Select.Item
          >
        {/if}
        {#each $AppStore.locations.filter((l) => l.parentLocationID === parent.id && !filterOutIds.includes(l.id)) as child}
          <Select.Item value={child.id}>{child.name}</Select.Item>
        {/each}
      {/each}
    {:else}
      {#each $AppStore.locations.filter((l) => l.parentLocationID && !filterOutIds.includes(l.id)) as child}
        <Select.Item value={child.id}>{child.name}</Select.Item>
      {/each}
    {/if}
  </Select.Content>
</Select.Root>
