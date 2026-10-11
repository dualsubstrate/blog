#!/usr/bin/env ruby
# frozen_string_literal: true

# Revise Psych-Fitness-28 English bundles (D4–D28): restrained **bold** in tldr;
# replace description with Sayings card teaser derived from title + tldr (+ fluff tone via petition line only if needed).
# Days D1–D3: skip (author-shaped hooks).

require "yaml"

ROOT = File.expand_path("..", __dir__)
GLOB = File.join(ROOT, "content/cognitive-memetics/psych-fitness-28/2026-05-*-day-*/index.md")

KEYWORDS = %w[
  psychological\ fitness mental\ fitness same\ human\ hardware tone\ at\ the\ top
  executive\ branch independent\ checks fitness\ checks public\ sector private\ sector
  mental\ energy nuclear\ war collective\ blindness civilian\ control reactive\ politics
  hard\ numbers modern\ governance checks\ and\ balances trust\ in\ government
  psychological\ fitness\ checks cognitive\ capacity accountability leadership institutions
  governance stability pressure stress crises crisis resilience taxpayers citizens
  mobilized warfare decisions signatures signature petition
].freeze

def strip_md_links(s)
  s.gsub(/\[([^\]]+)\]\([^)]+\)/, "\\1")
end

def plain(s)
  strip_md_links(s).gsub("**", "")
end

def emphasize_paragraph(para, max_bold: 3)
  return para if para.include?("**")

  n = 0
  out = para.dup
  KEYWORDS.each do |kw|
    break if n >= max_bold

    re = Regexp.new('\b' + Regexp.escape(kw) + '\b', Regexp::IGNORECASE)
    next unless out.match?(re)

    out.sub!(re) { |m| "**#{m}**" }
    n += 1
  end
  out
end

def emphasize_tldr(tldr)
  paras = tldr.strip.split(/\n\n+/).map(&:strip).reject(&:empty?)
  paras.map { |p| emphasize_paragraph(p, max_bold: 3) }.join("\n\n") + "\n"
end

def teaser_tail(tldr_text)
  paras = tldr_text.strip.split(/\n\n+/).map(&:strip).reject(&:empty?)
  # Use the FIRST paragraph for context, but keep it short
  first = paras.first.to_s.gsub(/\s+/, " ").strip
  first = first[0, 120] + (first.length > 120 ? "..." : "")
  
  # Add the PUNCH (last sentence of the last paragraph)
  last_para = paras.last.to_s.gsub(/\s+/, " ").strip
  sentences = last_para.split(/(?<=[.!?])\s+/)
  punch = sentences.last.to_s.strip

  t = "#{first} #{punch}"
  t = "#{t[0, 227]}..." if t.length > 230
  emphasize_paragraph(t, max_bold: 4)
end

def build_description(day, title, tldr_emphasized)
  tail = teaser_tail(tldr_emphasized)
  "**Day #{day} of 28**, **#{title}**: #{tail}"
end

def indent_tldr_yaml(new_body)
  new_body
    .rstrip
    .split(/\n\n+/, -1)
    .map(&:strip)
    .reject(&:empty?)
    .map { |para| para.each_line.map(&:chomp).map { |ln| "  #{ln}" }.join("\n") }
    .join("\n\n") + "\n"
end

def replace_tldr_block(text, new_body)
  indented = indent_tldr_yaml(new_body)
  re = /^tldr: \|\n([\s\S]*?)^fluff: \|\n/m
  unless text.match?(re)
    warn "skip tldr: pattern mismatch"
    return text
  end

  text.sub(re, "tldr: |\n#{indented}fluff: |\n")
end

def replace_description_line(text, new_desc)
  q =
    if new_desc.include?("'") && !new_desc.include?('"')
      '"' + new_desc.gsub('"', '\\"') + '"'
    else
      "'" + new_desc.gsub("'", "''") + "'"
    end
  # Single-line description only (do not use /m with .*)
  text.sub(/^description:[^\n]+\n/, "description: #{q}\n")
end

Dir.glob(GLOB).sort.each do |path|
  raw = File.read(path)
  parts = raw.split(/^---\s*$/m, -1)
  unless parts.size >= 3
    warn "skip #{path}: expected --- wrappers"
    next
  end

  fm = parts[1].strip
  body = parts[2]
  data = YAML.safe_load(fm)
  hc = data["heading_code"].to_s
  next if hc =~ /\AD[123]\z/

  m = hc.match(/\AD(\d{1,2})\z/)
  unless m
    warn "skip #{path}: heading_code #{hc.inspect}"
    next
  end

  day = m[1].to_i
  title = data["title"].to_s
  tldr = data["tldr"].to_s
  next if tldr.strip.empty?

  new_tldr = emphasize_tldr(tldr)
  new_desc = build_description(day, title, new_tldr)

  updated_fm = replace_tldr_block(fm, new_tldr)
  updated_fm = replace_description_line(updated_fm, new_desc)

  out = "---\n#{updated_fm}\n---#{body}"
  File.write(path, out)
  puts "updated #{File.basename(File.dirname(path))}"
end

puts "done"
