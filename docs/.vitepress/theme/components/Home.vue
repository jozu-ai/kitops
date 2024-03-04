<script setup lang="ts">
// @ts-expect-error
import { VueMarqueeSlider } from 'vue3-marquee-slider'
</script>

<template>
<div class="mt-32 md:mt-40 xl:mt-60 px-6 text-center">
  <p class="h4 !font-normal !text-off-white">Share and run your models anywhere</p>
  <h1 class="mt-4">Bridge the gap between ML and Application teams</h1>

  <a href="/docs/installation" class="kit-button mt-10 md:mt-14 xl:mt-22">Install</a>
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6 text-center max-w-[1152px] mx-auto">
  <h2>HoW does it WoRK<span class="font-heading font-extralight">?</span></h2>
  <div class="p1 mt-8 mx-8">
    Kit is an open source MLOps project that packages your model, datasets, code, and configuration so data scientists and developers can use their preferred tools while collaborating effortlessly.
  </div>

  <img src="/images/cli-screenshot.png" class="mt-22 mx-auto w-[1050px]" alt="screenshot of the kit cli tool" />
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6">
  <div class="kit-cards md:!grid-cols-2 w-fit mx-auto">
    <div class="kit-card max-w-[370px] flex flex-col">
      <h3>ModelKit</h3>

      <div class="flex-1 mt-8 space-y-4">
        <p class="p2">The ModelKit is an OCI compliant package that contains everything needed to integrate with a model, or deploy it to production. </p>
        <p class="p2">The ModelKit holds the serialized model, dataset, hyperparameters, input / output structure, and validation criteria. Kitfiles define a ModelKit in a modular and easy-to-understand way.</p>
      </div>

      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none" class="mt-10 md:mt-16 xl:mt-20">
        <path d="M6 16L18.0064 4H39.9956C41.1026 4 42 4.91062 42 5.9836V42.0164C42 43.112 41.1102 44 40.0132 44H7.9868C6.88952 44 6 43.1002 6 41.9864V16ZM20 8V18H10V40H38V8H20Z" class="fill-gold" />
      </svg>
    </div>
    <div class="kit-card max-w-[370px] flex flex-col">
      <h3 class="!text-cornflower">Kit cli</h3>

      <div class="flex-1 mt-8 space-y-4">
        <p class="p2">The Kit CLI is a command line interface (CLI) that performs actions on ModelKits.</p>
        <p class="p2">You can: build and version ModelKits; push or pull them from a model registry; run them locally with a RESTful API we generate for your model automatically, and deploy them to staging or production.</p>
      </div>

      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="32" viewBox="0 0 48 32" fill="none" class="mt-10 md:mt-16 xl:mt-20">
        <path d="M36.1238 25.0052H54.7501" class="stroke-cornflower" stroke-width="6.7732" stroke-linecap="square" stroke-linejoin="round"/>
        <path d="M0 5.25L15.8041 16.8209L0 26.6985" class="stroke-cornflower" stroke-width="6.7732" stroke-linecap="square" stroke-linejoin="round"/>
      </svg>
    </div>
  </div>
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6" id="whykitops">
  <h2 class="text-center">Why Kit<span class="font-heading font-extralight">?</span></h2>

  <div class="mt-10 md:mt-14 xl:mt-22 grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-x-4 md:gap-y-[4.5rem] xl:gap-22 max-w-[47.5rem] mx-auto">
    <div class="h4">Model handoffs are hard.</div>
    <div class="p2 space-y-4">
      <p>Moving a model from a Jupyter notebook to an ML tool or development server, then to a production server like Kubernetes is difficult because each tool uses its own packaging mechanism, and requires engineers to repackage the model multiple times. This slows down development and introduces risk.</p>
      <p>Kit is an open source MLOps project built to standardize packaging, reproduction, deployment, and tracking of AI / ML models, so it can be run anywhere, just like application code.</p>
      <p>Kit solves two big problems:</p>
    </div>

    <div class="h4 xs:mt-12">Collaboration–</div>
    <div class="p2 space-y-4">
      <p>By building ModelKits on industry standards, anyone (not just data scientists) can participate in the model development lifecycle whether they’re integrating models with their application, experimenting with them locally, or deploying them to production.</p>
      <p>ModelKits and Kitfiles work with the tools your team is already using, so you can use the same deployment pipelines and endpoints you’ve hardened with your application development process.</p>
    </div>

    <div class="h4 xs:mt-12">Model traceability and reproducibility–</div>
    <div class="p2 space-y-4">
      <p>Unlike Dockerfiles, Kitfiles are a modular package - pull just a part of the ModelKit, like the model or dataset, or pull the whole package with one simple command.</p>
      <p>Storing ModelKits in your organization’s registry provides a history of meaningful state changes for auditing. ModelKits are immutable so are perfect for a secure bill-of-materials (SBOM) initiative.</p>
    </div>
  </div>
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6" id="howdoesitwork">
  <h2 class="text-center">Get stArted</h2>

  <div class="kit-cards mt-22 min-h-[32.5rem]">
    <div class="kit-card flex flex-col">
      <div class="h4 font-bold !text-cornflower">1</div>
      <div class="mt-8 flex flex-col flex-1 justify-between">
        <p class="p2">Download and install Kit CLI.</p>
        <a href="/docs/installation" class="kit-button kit-button-salmon md:w-fit mt-6">Install the CLI</a>
      </div>
    </div>

    <div class="kit-card flex flex-col">
      <div class="h4 font-bold !text-cornflower">2</div>
      <div class="mt-8 flex flex-col flex-1 justify-between">
        <p class="p2">Create a simple manifest file called a Kitfile with your model, dataset and code. Then build and push the ModelKit to a registry for sharing.</p>
        <a href="/docs/installation" class="kit-button kit-button-cornflower md:w-fit mt-6">LEARN MORE</a>
      </div>
    </div>

    <div class="kit-card flex flex-col">
      <div class="h4 font-bold !text-cornflower">3</div>
      <div class="mt-8 flex flex-col flex-1 justify-between">
        <p class="p2">Pull the ModelKit into your pipeline, or use kit dev to start working with the model.</p>
        <a href="/docs/installation" class="kit-button md:w-fit mt-6">USE CASES</a>
      </div>
    </div>
  </div>
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6 text-center">
  <h2>WhAt’s suppoRted<span class="font-heading font-extralight">?</span></h2>
  <p class="p1 mt-8 mb-22">Kit was designed to work with the tools your team already uses.</p>

  <div class="space-y-12 relative marquee-gradients">
    <VueMarqueeSlider
      id="logos-row-1"
      :speed="30000"
      :space="75">
      <div class="flex justify-center items-center">
        <img src="/images/logos/jupyter@2x.png" alt="jupyter logo" class="opacity-65 hocus:opacity-100 max-h-10 md:max-h-12 xl:max-h-14">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/docker@2x.png" alt="docker logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/dvc@2x.png" alt="dvc logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/huggingface@2x.png" alt="huggingface logo" class="opacity-65 hocus:opacity-100 max-h-10 md:max-h-12 xl:max-h-14">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/github@2x.png" alt="github logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/gitlab@2x.png" alt="gitlab logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/jfrog@2x.png" alt="jfrog logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/azureml@2x.png" alt="azureml logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
    </VueMarqueeSlider>

    <VueMarqueeSlider
      id="logos-row-2"
      :speed="30000"
      :space="75"
      reverse>
      <div class="flex justify-center items-center">
        <img src="/images/logos/databricks@2x.png" alt="data bricks logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/datarobot@2x.png" alt="data robot logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/kubernetes@2x.png" alt="kubernetes logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/mlflow@2x.png" alt="mlflow logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/nvidia@2x.png" alt="nvidia logo" class="opacity-65 hocus:opacity-100 max-h-10 md:max-h-12 xl:max-h-14">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/openshift@2x.png" alt="openshift logo" class="opacity-65 hocus:opacity-100 max-h-10 md:max-h-12 xl:max-h-14">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/tensorflow@2x.png" alt="tensorflow logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
    </VueMarqueeSlider>

    <VueMarqueeSlider
      id="logos-row-3"
      :speed="30000"
      :space="75"
      :width="200">
      <div class="flex justify-center items-center">
        <img src="/images/logos/amazonsagemaker@2x.png" alt="amazon sage maker logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/circleci@2x.png" alt="circle ci logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/prefect@2x.png" alt="prefect logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/ray@2x.png" alt="ray logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/runai@2x.png" alt="runai logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/vertexai@2x.png" alt="vertex logo" class="opacity-65 hocus:opacity-100 max-h-8 md:max-h-10 xl:max-h-12">
      </div>
      <div class="flex justify-center items-center">
        <img src="/images/logos/weightsbiases@2x.png" alt="weights & biases logo" class="opacity-65 hocus:opacity-100 max-h-6 md:max-h-7 xl:max-h-8">
      </div>
    </VueMarqueeSlider>
  </div>

  <button class="kit-button mt-22">SEE FULL LIST</button>
</div>

<div class="mt-32 md:mt-40 xl:mt-60 px-6">
  <h2 class="text-center">How to Get inVolVeD<span class="font-heading font-extralight">?</span></h2>

  <div class="space-y-6 w-fit mx-auto mt-22">
    <a href="https://discord.gg/XzSmtPn3" class="border border-gray-02 p-8 md:px-14 md:py-10 flex justify-between gap-8 items-center hover:border-gold transition-colors">
      <div class="p1">Join the KitOps Discord</div>

      <svg xmlns="http://www.w3.org/2000/svg" width="43" height="32" viewBox="0 0 43 32" fill="none">
        <g clip-path="url(#clip0_8_909)">
          <path d="M36.476 2.67995C33.6728 1.40994 30.7134 0.508951 27.6736 0C27.2576 0.735868 26.8812 1.49298 26.546 2.26816C23.308 1.78531 20.0151 1.78531 16.777 2.26816C16.4416 1.49306 16.0653 0.735957 15.6495 0C12.6076 0.513249 9.64634 1.41638 6.84027 2.68659C1.26951 10.8427 -0.240636 18.7962 0.514437 26.6368C3.77681 29.022 7.42835 30.8361 11.3103 32C12.1844 30.8366 12.9579 29.6024 13.6225 28.3105C12.3601 27.8439 11.1417 27.2683 9.98138 26.5903C10.2868 26.3711 10.5854 26.1453 10.874 25.9261C14.2504 27.4974 17.9355 28.312 21.6666 28.312C25.3976 28.312 29.0827 27.4974 32.4591 25.9261C32.751 26.1619 33.0497 26.3877 33.3517 26.5903C32.1891 27.2694 30.9685 27.8461 29.7039 28.3138C30.3677 29.6052 31.1412 30.8383 32.0161 32C35.9014 30.8407 39.5557 29.0276 42.8187 26.6401C43.7046 17.5475 41.3052 9.66708 36.476 2.67995ZM14.5789 21.8149C12.4748 21.8149 10.7364 19.9253 10.7364 17.6007C10.7364 15.276 12.4144 13.3699 14.5722 13.3699C16.7301 13.3699 18.455 15.276 18.4181 17.6007C18.3811 19.9253 16.7233 21.8149 14.5789 21.8149ZM28.7542 21.8149C26.6467 21.8149 24.915 19.9253 24.915 17.6007C24.915 15.276 26.593 13.3699 28.7542 13.3699C30.9154 13.3699 32.6269 15.276 32.5899 17.6007C32.553 19.9253 30.8986 21.8149 28.7542 21.8149Z" fill="white"/>
        </g>
        <defs>
          <clipPath id="clip0_8_909">
            <rect width="42.6667" height="32" fill="white" transform="translate(0.333252)"/>
          </clipPath>
        </defs>
      </svg>
    </a>

    <a href="https://github.com/jozu-ai/kitops" class="border border-gray-02 p-8 md:px-14 md:py-10 flex justify-between gap-8 items-center hover:border-gold transition-colors">
      <div class="p1">Contribute to Kit</div>

      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
        <path d="M32.3432 21.9998L21.6152 11.2718L24.4436 8.44336L40 23.9998L24.4436 39.556L21.6152 36.7276L32.3432 25.9998H8V21.9998H32.3432Z" fill="#ECECEC"/>
      </svg>
    </a>

    <a href="https://github.com/jozu-ai/kitops" class="border border-gray-02 p-8 md:px-14 md:py-10 flex justify-between gap-8 items-center hover:border-gold transition-colors">
      <div class="p1">Star the repo on GitHub</div>

      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
        <g clip-path="url(#clip0_79_2469)">
          <path fill-rule="evenodd" clip-rule="evenodd" d="M23.9285 0C10.6967 0 0 11 0 24.6085C0 35.4865 6.85371 44.6945 16.3616 47.9535C17.5504 48.1985 17.9858 47.424 17.9858 46.7725C17.9858 46.202 17.9466 44.2465 17.9466 42.209C11.2903 43.676 9.90417 39.2755 9.90417 39.2755C8.83445 36.4235 7.24947 35.6905 7.24947 35.6905C5.07086 34.183 7.40816 34.183 7.40816 34.183C9.82482 34.346 11.0929 36.709 11.0929 36.709C13.2318 40.457 16.6785 39.398 18.0651 38.746C18.263 37.157 18.8973 36.057 19.5708 35.446C14.2619 34.8755 8.67625 32.757 8.67625 23.3045C8.67625 20.6155 9.62645 18.4155 11.1321 16.7045C10.8945 16.0935 10.0624 13.567 11.3701 10.1855C11.3701 10.1855 13.3905 9.5335 17.9461 12.7115C19.8965 12.1728 21.908 11.8988 23.9285 11.8965C25.9489 11.8965 28.0085 12.182 29.9104 12.7115C34.4665 9.5335 36.4869 10.1855 36.4869 10.1855C37.7946 13.567 36.962 16.0935 36.7244 16.7045C38.2697 18.4155 39.1807 20.6155 39.1807 23.3045C39.1807 32.757 33.5951 34.8345 28.2465 35.446C29.1184 36.22 29.8707 37.6865 29.8707 40.009C29.8707 43.309 29.8315 45.9575 29.8315 46.772C29.8315 47.424 30.2674 48.1985 31.4557 47.954C40.9636 44.694 47.8173 35.4865 47.8173 24.6085C47.8565 11 37.1207 0 23.9285 0Z" fill="white"/>
        </g>
        <defs>
          <clipPath id="clip0_79_2469">
            <rect width="48" height="48" fill="white"/>
          </clipPath>
        </defs>
      </svg>
    </a>
  </div>
</div>
</template>

<!-- Our custom home styles -->
<style scoped src="@theme/assets/css/home.css"></style>

<!-- The marquee lib css -->
<style src="../../../node_modules/vue3-marquee-slider/dist/style.css"></style>
