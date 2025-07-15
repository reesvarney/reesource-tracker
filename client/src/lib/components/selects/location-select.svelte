<script lang="ts">
  import LocationTree from "./LocationTree.svelte";
  import * as Select from "../ui/select";
  import { AppStore } from "$lib/components/app_store";
  import type { SampleLocation } from "$lib/components/location";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select a location";
  export let id: string = "location-select";
  export let required: boolean = false;
  export let options: { value: string; label: string }[] = [];
  export let onValueChange: (value: string) => void = () => {};
  export let filterMode: boolean = false;
  export let filterOutIds: string[] = [];

  let top_level_locations: SampleLocation[] = [];

  AppStore.subscribe((new_val) => {
    top_level_locations =
      new_val.locations.filter(
        (l: SampleLocation) =>
          l.parentLocationID === "" || l.parentLocationID === null
      ) ?? [];
  });
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
    {:else}
      <LocationTree
        locations={top_level_locations}
        {filterOutIds}
        level={0}
        {filterMode}
      />
    {/if}
  </Select.Content>
</Select.Root>
