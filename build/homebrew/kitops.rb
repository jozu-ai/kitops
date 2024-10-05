class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url 
      sha256 
    end
    on_intel do
      url 
      sha256 
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.16/kitops-linux-arm64.tar.gz"
      sha256 "b0a6857e1e58e161f63f1773b403b334b8e56e2d1b93b38147a2eacfca7bca2d"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.16/kitops-linux-x86_64.tar.gz"
        sha256 "6f14b82ed850f55e15fe8b6b27c39d614b3c74b3184b4914a8ad2455577911be"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.16/kitops-linux-i386.tar.gz"
        sha256 "3590f67562cd54033dbca40b67b9ad0bf93eb59d585ca69b90b14b0fb64c781e"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: "
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
