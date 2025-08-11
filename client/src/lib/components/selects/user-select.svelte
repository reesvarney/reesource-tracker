<script lang="ts">
  import * as Select from "../ui/select";
  import { AppStore } from "$lib/components/app_store";
  import type { User } from "$lib/components/user";
  export let bindValue: string;
  export let disabled: boolean = false;
  export let placeholder: string = "Select a user";
  export let id: string = "user-select";
  export let required: boolean = false;
  export let options: { value: string; label: string }[] = [];
  export let onValueChange: (value: string) => void = () => {};
  export let filterOutIds: string[] = [];
  export let filterMode: boolean = false;
  let users: User[] = [];

  AppStore.subscribe((new_val) => {
    users = new_val.users ?? [];
  });

  function findUserById(id: string): User | null {
    return users.find((user) => user.id === id) || null;
  }
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
      {:else}
        {findUserById(bindValue)?.name}
      {/if}
    {:else}
      {placeholder}
    {/if}
  </Select.Trigger>
  <Select.Content class="w-full">
    {#if filterMode}
      <Select.Item value="">Any</Select.Item>
    {:else}
      <Select.Item value="">None</Select.Item>
    {/if}
    {#if options.length > 0}
      {#each options.filter((option) => !filterOutIds.includes(option.value) && option.label !== "") as option}
        <Select.Item value={option.value}>{option.label}</Select.Item>
      {/each}
    {:else}
      {#each users.filter((user) => !filterOutIds.includes(user.id) && user.name !== "") as user}
        <Select.Item value={user.id}>{user.name}</Select.Item>
      {/each}
    {/if}
  </Select.Content>
</Select.Root>
