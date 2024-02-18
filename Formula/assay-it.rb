# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class AssayIt < Formula
  desc "Confirm Quality and Eliminate Risk by Testing Microservices in Production."
  homepage "https://assay.it"
  version "1.2.12"
  license "MIT"

  depends_on "go"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.12/assay-it_1.2.12_darwin_amd64"
      sha256 "b5554c378d6c10c05b85315bc91a6c6f3749b2b07acbb5e29a3766da6455825b"

      def install
        bin.install "assay-it_1.2.12_darwin_amd64" => "assay-it"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.12/assay-it_1.2.12_darwin_arm64"
      sha256 "a58845af9ee40508fe09104e003e2d2ef2b864714bd8870f23ac36dfce6d9d6f"

      def install
        bin.install "assay-it_1.2.12_darwin_arm64" => "assay-it"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.12/assay-it_1.2.12_linux_arm64"
      sha256 "079edff6b9fc6d177e58539ec7c3ad39d2bf7090fff185d4e6a00490fb350e1c"

      def install
        bin.install "assay-it_1.2.12_linux_arm64" => "assay-it"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.12/assay-it_1.2.12_linux_amd64"
      sha256 "9428af8f39baf58a0d89ad5a6ffdbbeff2b48edac907797d662d1a94d4ac4009"

      def install
        bin.install "assay-it_1.2.12_linux_amd64" => "assay-it"
      end
    end
  end

  test do
    system "#{bin}/assay-it -v"
  end
end
