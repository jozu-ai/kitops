class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-darwin-arm64.tar.gz"
      sha256 "df42727a4b330427888cd6f1a05367105173a6b8b3add0da58cef441423f1aaa"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-darwin-x86_64.tar.gz"
      sha256 "5120c214de5c2c328876eacea0bf88b3e61ae9fb98adf263bfd48ac9627e6764"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-arm64.tar.gz"
      sha256 "2862f7c28bfbc12d3e8e58b01194892bd4df8016a749b33b03fdeda2aac710a6"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-x86_64.tar.gz"
        sha256 "20605a8bacc0aec96011b3d9114b503c5a348165e0603e397820284e15184483"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-i386.tar.gz"
        sha256 "93e704ebb7aeba96c7070f754e02023b000fe5f74c4b1d5e5bd6fda14307311e"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: 0.4.4"
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
