# What is KitOps?

KitOps is an innovative open-source project designed to enhance collaboration among data scientists, application developers, and SREs working on integrating or managing self-hosted AI/ML models.

## What's Inside KitOps?

### 🎁 ModelKit

At the heart of KitOps is the ModelKit, an OCI-compliant packaging format that enables the seamless sharing of all necessary artifacts involved in the AI/ML model lifecycle. This includes datasets, code, configurations, and the models themselves. By standardizing the way these components are packaged, ModelKit facilitates a more streamlined and collaborative development process that is compatible with nearly any tool. You can even [deploy ModelKits to containers or Kubernetes](./deploy.md).

### 📄 Kitfile

Complementing the ModelKit is the Kitfile, a YAML-based configuration file that simplifies the sharing of model, dataset, and code configurations. The Kitfile is designed with both ease of use and security in mind, ensuring that configurations can be efficiently packaged and shared without compromising on safety or governance.

### 🖥️ Kit CLI

Bringing everything together is the Kit Command Line Interface (CLI). The Kit CLI is a powerful tool that enables users to create, manage, run, and deploy ModelKits using Kitfiles. Whether you are packaging a new model for development or deploying an existing model into production, the Kit CLI provides the necessary commands and functionalities to streamline your workflow.

## How KitOps is Used

KitOps is a key element in a platform engineering solution for AI/ML projects.

[See how security-conscious organization are using ModelKits](./use-cases.md) with their existing tools to develop AI/ML projects faster and safer than ever before.

## The Goal of KitOps

The primary goal of KitOps is to become an open, vendor-neutral standard that simplifies and secures the packaging and versioning of AI/ML projects. In the same way that PDFs have helped people share documents, images, and diagrams between tools, KitOps makes it easy for teams to use the tools they prefer, but share the results safely and securely.

KitOps drives greater speed, security, and collaboration for teams working with models.

### 👩‍💻 For application developers

KitOps clears the path to use AI/ML with your existing tools and applications. No need to be an AI/ML expert, KitOps lets you concentrate on integrating AI/ML models into your applications, while Kit handles the packaging and sharing.

[Get Started](./get-started.md).

### 👷 For DevOps teams

ModelKits fit into your existing processes and the Kit CLI lets you pack or unpack ModelKit artifacts in the pipelines and automation you have proven over the last decade.

[Build a better golden path for AI/ML projects](./use-cases.md).
[Get Started](./get-started.md).


### 👩‍🔬 For data scientists

KitOps enables you to innovate in AI/ML without the usual infrastructure distractions. It simplifies dataset and model management and sharing, fostering closer collaboration with developers. With KitOps, you can spend more time experimenting and less time grappling with traditional software development tools.

[See how to use KitOps with Jupyter Notebooks](https://www.youtube.com/watch?v=OQPp7QEvk7Q).
[Get Started](./get-started.md).

## Benefits of KitOps

KitOps is not just another tool; it's a comprehensive CLI and packaging system specifically designed for the AI/ML workflow. It acknowledges the nuanced needs of AI/ML projects, such as:

### 📊 Management of Unstructured Datasets

AI/ML projects often deal with large, unstructured datasets, such as images, videos, and audio files. KitOps simplifies the versioning and sharing of these datasets, making them as manageable as traditional code.

### 🤝 Synchronized Data and Code Versioning

One of the core strengths of KitOps is its ability to keep data and code versions in sync. This crucial feature solves the reproducibility issues that frequently arise in AI/ML development, ensuring consistency and reliability across project stages.

### 🚀 Deployment Ready

Designed with a focus on deployment, ModelKits package assets in standard formats so you can depoly them as [containers or to Kubernetes](./deploy.md). They're also [compatible with nearly any tool](./modelkit/compatibility.md) - helping you get your model to production faster and more efficiently.

### 🏭 Standards-Based Approach

KitOps champions openness and interoperability through its core components, ensuring seamless integration into your existing workflows:

ModelKits are designed as OCI (Open Container Initiative) artifacts, making them fully compatible with the Docker image registries and other OCI-compliant storage solutions you already use. This compatibility allows for an easy and familiar integration process. By adhering to widely accepted standards, KitOps ensures you're not tied to a single vendor or platform. This flexibility gives you the freedom to choose the best tools and services for your needs without being restricted by proprietary formats.

Kitfiles leverage the simplicity and ubiquity of YAML for configuration, offering an accessible and straightforward way to specify the details of your AI/ML projects.

The Kit CLI is an open-source tool, developed and supported by a community passionate about advancing AI/ML collaboration. Its open-source nature not only fosters innovation and continuous improvement but also allows you to customize and extend its capabilities to meet your unique project requirements.
