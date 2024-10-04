class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-darwin-arm64.tar.gz"
      sha256 "05fb4154c4db321d0498767c30d8d54e1cd4c028e7500bd14446f72e619b2f23"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-darwin-x86_64.tar.gz"
      sha256 "7f8ad1c5a33979a8ccef0492bc37c19974f15a3dc02d6c3f9f3e53bd1b69153d"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-arm64.tar.gz"
      sha256 "6a5cebafee8b452615a3cfdb9825bef0cddaa3c0eb7fcbf9cda7e8e94a6f542f"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-x86_64.tar.gz"
        sha256 "efc429106df85e580656aabc27d3c442362fd98d0351e623ce7880aa0b5c05aa"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.4/kitops-linux-i386.tar.gz"
        sha256 "3fead337eff880bad9ee898de06ed2f399ea6f38fb1120141d538d8a66bd38a4"
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
