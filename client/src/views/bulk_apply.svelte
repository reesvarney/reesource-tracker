<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import QRScanner from "$lib/components/qr_scanner/qr_scanner.svelte";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import StateSelect from "$lib/components/selects/state-select.svelte";
  import { toast } from "svelte-sonner";
  import * as InputOTP from "$lib/components/ui/input-otp";
  let scannedIds: string[] = $state([]);
  let selectedProduct = $state("");
  let selectedLocation = $state("");
  let sampleState = $state("");

  let { active = $bindable(false) } = $props();
  // QR scan handler
  function handleQRScan(text: string) {
    let sampleId: string | null = null;
    try {
      // Try to parse as URL and extract sample_id param
      const url = new URL(text);

      const param = url.searchParams.get("sample_id");
      if (param && /^[0-9A-Z]{2}-[0-9A-Z]{2}-[0-9A-Z]{2}$/i.test(param)) {
        sampleId = param.toUpperCase();
      }
    } catch (e) {
      // Not a valid URL, fallback to direct code
      if (/^[0-9A-Z]{2}-[0-9A-Z]{2}-[0-9A-Z]{2}$/i.test(text)) {
        sampleId = text.toUpperCase();
      }
    }
    if (sampleId) {
      // Only add if not already present (case-insensitive)
      const exists = scannedIds.some(
        (id) => id.toUpperCase() === sampleId!.toUpperCase()
      );
      if (!exists) {
        toast.success(`Scanned sample ID: ${sampleId}`);
        scannedIds = [...scannedIds, sampleId];
      }
    }
  }

  let selectedVideoInput: string = $state("");

  // Manual sample code entry
  let manualSample: string = $state("");
  let manualSampleError: string = $state("");

  $effect(() => {
    if (manualSample.length === 6) {
      const formatted =
        `${manualSample.slice(0, 2)}-${manualSample.slice(2, 4)}-${manualSample.slice(4, 6)}`.toUpperCase();
      if (/^[0-9A-Z]{2}-[0-9A-Z]{2}-[0-9A-Z]{2}$/.test(formatted)) {
        // Only add if not already present (case-insensitive)
        const exists = scannedIds.some((id) => id.toUpperCase() === formatted);
        if (!exists) {
          scannedIds = [...scannedIds, formatted];
          toast.success(`Added sample ID: ${formatted}`);
          manualSampleError = "";
        } else {
          manualSampleError = "Sample ID already added.";
        }
      } else {
        manualSampleError = "Invalid sample code format.";
      }
      manualSample = "";
    }
  });

  let updateProduct = $state(true);
  let updateLocation = $state(true);
  let updateState = $state(true);
  let modNames: string[] = $state([]);
  let modInput = $state("");
  let modError = $state("");
  let modLoading = $state(false);

  async function applyChanges() {
    if (
      !updateProduct &&
      !updateLocation &&
      !updateState &&
      modNames.length === 0
    ) {
      toast.error(
        "Please enable at least one field to update or add at least one mod."
      );
      return;
    }
    if (scannedIds.length === 0) {
      toast.error("No samples scanned.");
      return;
    }
    let errors = 0;
    let modErrors = 0;
    for (const id of scannedIds) {
      let productId = selectedProduct;
      let locationId = selectedLocation;
      let state = sampleState;
      let sampleFetched = false;
      // If a field is not enabled, fetch the current value for this sample
      if (!updateProduct || !updateLocation || !updateState) {
        const res = await fetch(`/api/sample/${id}`);
        if (res.ok) {
          sampleFetched = true;
          const data = await res.json();
          if (!updateProduct) productId = data.sample.ProductID || "";
          if (!updateLocation) locationId = data.sample.LocationID || "";
          if (!updateState) state = data.sample.State || "";
        }
      }
      // If sample not found and state is not configured, fallback to 'available'
      if (!sampleFetched && !updateState) {
        state = "available";
      }
      // Apply main changes
      if (updateProduct || updateLocation || updateState) {
        const form = new FormData();
        form.append("product_id", productId);
        form.append("location_id", locationId);
        form.append("state", state);
        const res = await fetch(`/api/sample/${id}`, {
          method: "POST",
          body: form,
        });
        if (!res.ok) errors++;
      }
      // Apply all mods
      for (const mod of modNames) {
        if (mod.trim()) {
          const res = await fetch(`/api/sample/${id}/mods`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name: mod }),
          });
          if (!res.ok) modErrors++;
        }
      }
    }
    let msg = [];
    if ((updateProduct || updateLocation || updateState) && errors === 0) {
      msg.push("Changes applied to all scanned samples.");
    } else if ((updateProduct || updateLocation || updateState) && errors > 0) {
      msg.push(`Some updates failed (${errors} errors).`);
    }
    if (modNames.length > 0) {
      if (modErrors === 0) {
        msg.push("All mods added to all scanned samples.");
        modNames = [];
      } else {
        msg.push(`Some mods failed to add (${modErrors} errors).`);
      }
    }
    if (msg.length > 0) {
      toast(msg.join("\n"));
    }
  }

  function clearScanned() {
    scannedIds = [];
    toast("Cleared all scanned sample IDs.");
  }
</script>

<div class="space-y-4">
  <div class="flex flex-col items-center shrink">
    <label for="qr-reader-bulk" class="mb-6">Scan a QR code</label>
    <QRScanner
      containerId="qr-reader-bulk"
      bind:selectedVideoInput
      onQrCodeScan={handleQRScan}
      autoStart={active}
    />
  </div>
  <div class="self-center mt-6 flex flex-col items-center gap-6">
    <label for="manual-id-input">Or manually enter the sample ID</label>
    <InputOTP.Root
      maxlength={6}
      bind:value={manualSample}
      id="manual-id-input"
      disabled={modLoading}
    >
      {#snippet children({ cells })}
        <InputOTP.Group>
          {#each cells.slice(0, 2) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
        <InputOTP.Separator />
        <InputOTP.Group>
          {#each cells.slice(2, 4) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
        <InputOTP.Separator />
        <InputOTP.Group>
          {#each cells.slice(4, 6) as cell}
            <InputOTP.Slot {cell} />
          {/each}
        </InputOTP.Group>
      {/snippet}
    </InputOTP.Root>
    {#if manualSampleError}
      <span class="text-red-500 text-sm">{manualSampleError}</span>
    {/if}
  </div>

  <div>
    <label class="font-bold flex items-center gap-2">
      <input type="checkbox" bind:checked={updateProduct} />
      Product
    </label>
    <ProductSelect bind:bindValue={selectedProduct} disabled={!updateProduct} />
  </div>

  <div>
    <label class="font-bold flex items-center gap-2">
      <input type="checkbox" bind:checked={updateLocation} />
      Location
    </label>
    <LocationSelect
      bind:bindValue={selectedLocation}
      disabled={!updateLocation}
    />
  </div>

  <div>
    <label class="font-bold flex items-center gap-2">
      <input type="checkbox" bind:checked={updateState} />
      State
    </label>
    <StateSelect bind:bindValue={sampleState} disabled={!updateState} />
  </div>

  <div class="flex flex-col gap-2">
    <div class="flex gap-2 items-center">
      <input
        type="text"
        class="border rounded px-2 py-1 flex-1"
        placeholder="Add mod (optional, applies with changes)"
        bind:value={modInput}
        onkeydown={(e) => {
          if (e.key === "Enter" && modInput.trim()) {
            e.preventDefault();
            if (!modNames.includes(modInput.trim())) {
              modNames = [...modNames, modInput.trim()];
            }
            modInput = "";
          }
        }}
        disabled={modLoading}
      />
      <Button
        onclick={() => {
          if (modInput.trim() && !modNames.includes(modInput.trim())) {
            modNames = [...modNames, modInput.trim()];
            modInput = "";
          }
        }}
        disabled={modLoading || !modInput.trim()}>Add</Button
      >
      {#if modError}
        <span class="text-red-500 text-sm">{modError}</span>
      {/if}
    </div>
    {#if modNames.length > 0}
      <div class="flex flex-wrap gap-2 mb-2">
        {#each modNames as mod, i}
          <span class="bg-gray-200 rounded px-2 py-1 flex items-center gap-1">
            {mod}
            <button
              type="button"
              class="ml-1 text-red-500 hover:text-red-700"
              onclick={() => {
                modNames = modNames.filter((_, idx) => idx !== i);
              }}
              aria-label="Remove mod">&times;</button
            >
          </span>
        {/each}
      </div>
    {/if}
    <Button
      onclick={applyChanges}
      disabled={scannedIds.length === 0 || modLoading}>Apply Changes</Button
    >
    <Button
      variant="secondary"
      onclick={clearScanned}
      disabled={scannedIds.length === 0 || modLoading}>Clear Scanned</Button
    >
  </div>

  <div class="mt-4">
    <h2 class="font-bold">Scanned Sample IDs</h2>
    {#if scannedIds.length === 0}
      <div class="text-muted-foreground">No samples scanned yet.</div>
    {:else}
      <ul class="list-disc ml-6">
        {#each scannedIds as id}
          <li>{id}</li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

<style>
  #qr-reader-bulk {
    width: 100%;
    max-width: 400px;
    margin: auto;
  }
</style>
