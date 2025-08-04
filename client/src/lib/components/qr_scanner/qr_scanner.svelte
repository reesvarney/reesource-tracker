<script lang="ts">
  import { onMount, onDestroy, createEventDispatcher } from "svelte";
  import QrScanner from "qr-scanner";

  let videoInputs: MediaDeviceInfo[] = $state([]);
  let {
    selectedVideoInput = $bindable(""),
    containerId = "qr-reader",
    autoStart = $bindable(false),
    flipFrontCamera = $bindable(true),
    onQrCodeScan = $bindable((text: string) => {}),
  } = $props();
  let videoInputError: string = $state("");

  let qrScanner: QrScanner | null = null;
  let videoElement: HTMLVideoElement | null = null;
  let scanning = false;

  async function startScanner() {
    if (scanning) return;
    scanning = true;
    const qrContainer = document.getElementById(containerId);
    if (!qrContainer) return;
    qrContainer.innerHTML = "";
    videoElement = document.createElement("video");
    videoElement.className = "max-h-[20vh]";
    videoElement.setAttribute("autoplay", "true");
    videoElement.setAttribute("playsinline", "true");
    videoElement.style.width = "100%";

    // Flip preview if front-facing camera is selected
    if (flipFrontCamera && videoInputs && selectedVideoInput) {
      const selectedDevice = videoInputs.find(
        (d) => d.deviceId === selectedVideoInput
      );
      if (
        selectedDevice &&
        /front|user|integrated/i.test(selectedDevice.label)
      ) {
        videoElement.style.transform = "scaleX(-1)";
      } else {
        videoElement.style.transform = "";
      }
    }

    qrContainer.appendChild(videoElement);

    // Stop any previous scanner
    if (qrScanner) {
      qrScanner.destroy();
      qrScanner = null;
    }

    qrScanner = new QrScanner(
      videoElement,
      (result) => {
        if (result) {
          onQrCodeScan(result.data);
        }
      },
      {
        preferredCamera: selectedVideoInput || undefined,
        highlightScanRegion: false,
        highlightCodeOutline: true,
        maxScansPerSecond: 30,
        calculateScanRegion(video) {
          return {
            x: 0,
            y: 0,
            width: video.videoWidth,
            height: video.videoHeight,
          };
        },
        // Use high-res constraints for the stream
      }
    );
    await qrScanner.start();
  }

  function stopScanner() {
    scanning = false;
    if (qrScanner) {
      qrScanner.destroy();
      qrScanner = null;
    }
    if (videoElement && videoElement.srcObject) {
      const tracks = (videoElement.srcObject as MediaStream).getTracks();
      tracks.forEach((track) => track.stop());
      videoElement.srcObject = null;
    }
    videoElement = null;
    const qrContainer = document.getElementById(containerId);
    if (qrContainer) qrContainer.innerHTML = "";
  }

  $effect(() => {
    if (autoStart && selectedVideoInput) {
      restartScanner();
    } else if (!autoStart) {
      stopScanner();
    }
  });

  async function restartScanner() {
    stopScanner();
    await startScanner();
  }

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

  onMount(async () => {
    await enumerateVideoInputs();
    if (autoStart && selectedVideoInput) startScanner();
  });
  onDestroy(() => {
    stopScanner();
  });
</script>

<div id={containerId} class="border p-4"></div>
<div class="mt-2">
  <label class="block font-bold" for="camera-select">Camera</label>
  {#if videoInputError}
    <div class="text-destructive">{videoInputError}</div>
  {:else}
    <select
      id="camera-select"
      bind:value={selectedVideoInput}
      class="border rounded px-2 py-1"
    >
      {#each videoInputs as device}
        <option value={device.deviceId}>
          {device.label || `Camera ${device.deviceId}`}
        </option>
      {/each}
    </select>
  {/if}
</div>

<style>
  .qr-reader {
    width: 100%;
    max-width: 400px;
    margin: auto;
  }
</style>
