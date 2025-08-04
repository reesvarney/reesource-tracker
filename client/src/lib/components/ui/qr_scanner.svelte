<script lang="ts"></script>
  import { onMount, onDestroy, createEventDispatcher } from "svelte";
  import { BrowserQRCodeReader } from "@zxing/browser";

  export let videoInputs: MediaDeviceInfo[] = [];
  export let selectedVideoInput: string = "";
  export let videoInputError: string = "";
  export let containerId: string = "qr-reader";
  export let autoStart: boolean = true;
  export let flipFrontCamera: boolean = true;

  // Event: fires with the scanned text (string)
  const dispatch = createEventDispatcher();

  let codeReader: BrowserQRCodeReader | null = null;
  let videoElement: HTMLVideoElement | null = null;
  let scanning = false;

  async function startScanner() {
    if (scanning) return;
    codeReader = new BrowserQRCodeReader();
    scanning = true;
    const qrContainer = document.getElementById(containerId);
    if (!qrContainer) return;
    qrContainer.innerHTML = "";
    videoElement = document.createElement("video");
    videoElement.setAttribute("autoplay", "true");
    videoElement.setAttribute("playsinline", "true");
    videoElement.style.width = "100%";

    // Flip preview if front-facing camera is selected
    if (flipFrontCamera && videoInputs && selectedVideoInput) {
      const selectedDevice = videoInputs.find(
        (d) => d.deviceId === selectedVideoInput
      );
      if (selectedDevice && /front|user|integrated/i.test(selectedDevice.label)) {
        videoElement.style.transform = "scaleX(-1)";
      } else {
        videoElement.style.transform = "";
      }
    }

    qrContainer.appendChild(videoElement);

    codeReader.decodeFromVideoDevice(
      selectedVideoInput || undefined,
      videoElement,
      (result, err, controls) => {
        if (result && result.getText) {
          dispatch("scan", { text: result.getText() });
        }
        // Ignore errors
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
    const qrContainer = document.getElementById(containerId);
    if (qrContainer) qrContainer.innerHTML = "";
  }

  $: if (autoStart && selectedVideoInput) {
    restartScanner();
  }

  async function restartScanner() {
    stopScanner();
    await startScanner();
  }

  onMount(() => {
    if (autoStart && selectedVideoInput) startScanner();
  });
  onDestroy(() => {
    stopScanner();
  });
</script>

<div id={containerId} class="border p-4"></div>

<style>
  .qr-reader {
    width: 100%;
    max-width: 400px;
    margin: auto;
  }
</style>
