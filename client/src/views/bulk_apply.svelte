<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Select from "$lib/components/ui/select";
  import { BrowserQRCodeReader } from "@zxing/browser";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import StateSelect from "$lib/components/selects/state-select.svelte";
  let scannedIds: string[] = $state([]);
  let selectedProduct = $state("");
  let selectedLocation = $state("");
  let sampleState = $state("");
  let codeReader: BrowserQRCodeReader | null = null;
  let videoElement: HTMLVideoElement | null = null;
  let scanning = false;

  let videoInputs: MediaDeviceInfo[] = $state([]);
  let selectedVideoInput: string = $state("");
  let videoInputError: string = $state("");

  let updateProduct = $state(true);
  let updateLocation = $state(true);
  let updateState = $state(true);

  async function enumerateVideoInputs() {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      videoInputs = devices.filter((d) => d.kind === "videoinput");
      if (videoInputs.length > 0) {
        if (
          !selectedVideoInput ||
          !videoInputs.find((d) => d.deviceId === selectedVideoInput)
        ) {
          selectedVideoInput = videoInputs[0].deviceId;
        }
        videoInputError = "";
      } else {
        videoInputError = "No video input devices found.";
        selectedVideoInput = "";
      }
    } catch (e) {
      videoInputError = "Unable to enumerate video devices.";
      videoInputs = [];
      selectedVideoInput = "";
    }
  }

  async function startScanner() {
    if (scanning) return;
    codeReader = new BrowserQRCodeReader();
    scanning = true;
    const qrContainer = document.getElementById("qr-reader-bulk");
    if (!qrContainer) return;
    qrContainer.innerHTML = "";
    videoElement = document.createElement("video");
    videoElement.setAttribute("autoplay", "true");
    videoElement.setAttribute("playsinline", "true");
    videoElement.style.width = "100%";
    qrContainer.appendChild(videoElement);

    // Ensure the video stream is attached to the video element for preview
    codeReader.decodeFromVideoDevice(
      selectedVideoInput || undefined,
      videoElement,
      (result, err, controls) => {
        if (result && result.getText) {
          const text = result.getText();
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
              (id) => id.toUpperCase() === sampleId.toUpperCase()
            );
            if (!exists) {
              scannedIds = [...scannedIds, sampleId];
            }
          }
        }
        // Ignore errors (err) for now
      }
    );
  }

  $effect(() => {
    if (selectedVideoInput) {
      restartScanner();
    }
  });

  async function restartScanner() {
    scanning = false;
    await startScanner();
  }

  // Mods quick add/remove for all scanned samples
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
      alert(
        "Please enable at least one field to update or add at least one mod."
      );
      return;
    }
    if (scannedIds.length === 0) {
      alert("No samples scanned.");
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
      alert(msg.join("\n"));
    }
  }

  function clearScanned() {
    scannedIds = [];
  }

  onMount(async () => {
    await enumerateVideoInputs();
    await startScanner();
  });
  function stopScanner() {
    scanning = false;
    if (codeReader) {
      codeReader = null;
    }
    if (videoElement && videoElement.srcObject) {
      const tracks = (videoElement.srcObject as MediaStream).getTracks();
      tracks.forEach((track) => track.stop());
      videoElement.srcObject = null;
    }
    videoElement = null;
    const qrContainer = document.getElementById("qr-reader-bulk");
    if (qrContainer) qrContainer.innerHTML = "";
  }

  onDestroy(() => {
    stopScanner();
  });
</script>

<div class="space-y-4">
  <div>
    <label class="block font-bold" for="camera-select">Camera</label>
    {#if videoInputError}
      <div class="text-destructive">{videoInputError}</div>
    {:else}
      <Select.Root bind:value={selectedVideoInput} type="single">
        <Select.Trigger class="w-full" id="camera-select">
          {selectedVideoInput
            ? videoInputs.find((d) => d.deviceId === selectedVideoInput)
                ?.label || "Camera"
            : "Select camera"}
        </Select.Trigger>
        <Select.Content class="w-full">
          {#each videoInputs as device}
            <Select.Item value={device.deviceId}
              >{device.label || `Camera ${device.deviceId}`}</Select.Item
            >
          {/each}
        </Select.Content>
      </Select.Root>
    {/if}
  </div>

  <div id="qr-reader-bulk" class="border p-4"></div>

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
