<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { BrowserQRCodeReader } from "@zxing/browser";
  import * as InputOTP from "$lib/components/ui/input-otp";

  let codeReader: BrowserQRCodeReader | null = null;
  let videoElement: HTMLVideoElement | null = null;
  let scanning = false;
  let videoInputs: MediaDeviceInfo[] = $state([]);
  let selectedVideoInput: string = $state("");
  let videoInputError: string = $state("");
  let visible = false;
  let rootEl: HTMLElement | null = null;
  let observer: IntersectionObserver | null = null;
  let new_sample: string = $state("");

  $effect(() => {
    if (new_sample.length == 6) {
      window.location.assign(
        `/app?sample_id=${new_sample.slice(0, 2)}-${new_sample.slice(2, 4)}-${new_sample.slice(4, 6)}`
      );
    }
  });
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
    if (scanning || !visible) return;
    codeReader = new BrowserQRCodeReader();
    scanning = true;
    const qrContainer = document.getElementById("qr-reader-find");
    if (!qrContainer) return;
    qrContainer.innerHTML = "";
    videoElement = document.createElement("video");
    videoElement.setAttribute("autoplay", "true");
    videoElement.setAttribute("playsinline", "true");
    videoElement.style.width = "100%";
    qrContainer.appendChild(videoElement);

    await codeReader.decodeFromVideoDevice(
      selectedVideoInput || undefined,
      videoElement,
      (result, err, controls) => {
        if (result && result.getText) {
          const text = result.getText();
          try {
            const url = new URL(text);
            // Only open if same host
            if (url.host === window.location.host) {
              window.location.href = url.href;
            } else {
              alert("Scanned QR code is not for this host.");
            }
          } catch (e) {
            alert("Scanned QR code is not a valid URL.");
          }
        }
        // Ignore errors (err) for now
      }
    );
  }

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
    const qrContainer = document.getElementById("qr-reader-find");
    if (qrContainer) qrContainer.innerHTML = "";
  }

  $effect(() => {
    if (selectedVideoInput && visible) {
      restartScanner();
    } else if (!visible) {
      stopScanner();
    }
  });

  async function restartScanner() {
    stopScanner();
    await startScanner();
  }

  onMount(async () => {
    await enumerateVideoInputs();
    observer = new IntersectionObserver(
      (entries) => {
        visible = entries[0]?.isIntersecting ?? false;
        if (visible) {
          startScanner();
        } else {
          stopScanner();
        }
      },
      { threshold: 0.1 }
    );
    if (rootEl) observer.observe(rootEl);
  });
  onDestroy(() => {
    scanning = false;
    if (observer && rootEl) observer.unobserve(rootEl);
    observer = null;
    stopScanner();
  });
</script>

<div class="space-y-4 flex flex-col" bind:this={rootEl}>
  <div class="flex flex-col items-center">
    <label for="qr-reader" class="mb-6">Scan a QR code</label>
    <div id="qr-reader-find" class="border p-4"></div>

    <label class="block font-bold mt-6" for="camera-select">Camera</label>
    {#if videoInputError}
      <div class="text-destructive">{videoInputError}</div>
    {:else}
      <select
        id="camera-select"
        bind:value={selectedVideoInput}
        class="border rounded px-2 py-1"
      >
        {#each videoInputs as device}
          <option value={device.deviceId}
            >{device.label || `Camera ${device.deviceId}`}</option
          >
        {/each}
      </select>
    {/if}
  </div>
  <div class="self-center mt-12 flex flex-col items-center gap-6">
    <label for="id-input">Or manually enter the sample ID</label>
    <InputOTP.Root maxlength={6} bind:value={new_sample} id="id-input">
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
  </div>
</div>

<style>
  #qr-reader-find {
    width: 100%;
    max-width: 400px;
    margin: auto;
  }
</style>
