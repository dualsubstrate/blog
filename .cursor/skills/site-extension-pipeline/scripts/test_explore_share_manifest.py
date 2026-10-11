#!/usr/bin/env python3
"""Unit tests for explore_share_manifest preflight."""

from __future__ import annotations

import unittest
from datetime import timedelta

from explore_share_manifest import (
    SHARE_SETTING_ANYONE,
    empty_manifest,
    isoformat_z,
    prepare_rows_from_packet,
    utc_now,
    validate_manifest_for_urls,
    validate_row,
)


class ExploreShareManifestTests(unittest.TestCase):
    def _row(self, **overrides):
        now = utc_now()
        base = {
            "candidate_id": "blame-blindspot",
            "url": "https://www.perplexity.ai/search/bacb0d15-2888-45ab-ab32-b30aae743709",
            "thread_id": "bacb0d15-2888-45ab-ab32-b30aae743709",
            "share_setting": SHARE_SETTING_ANYONE,
            "share_confirmed_at": isoformat_z(now),
            "cold_verified": True,
            "cold_verified_at": isoformat_z(now),
            "cold_verification_method": "operator_incognito",
        }
        base.update(overrides)
        return base

    def test_prepare_rows_from_packet(self):
        packet = {
            "candidates": [
                {
                    "id": "a",
                    "url": "https://www.perplexity.ai/search/11111111-1111-4111-8111-111111111111",
                },
                {"id": "b", "url": ""},
            ]
        }
        rows = prepare_rows_from_packet(packet)
        self.assertEqual(len(rows), 1)
        self.assertEqual(rows[0]["candidate_id"], "a")

    def test_validate_row_ok(self):
        errors = validate_row(self._row())
        self.assertEqual(errors, [])

    def test_validate_row_missing_cold(self):
        errors = validate_row(self._row(cold_verified=False, cold_verified_at=None))
        self.assertTrue(any("cold_verified" in e for e in errors))

    def test_validate_row_stale(self):
        old = utc_now() - timedelta(days=40)
        errors = validate_row(
            self._row(cold_verified_at=isoformat_z(old)),
            now=utc_now(),
            max_cold_age_days=30,
        )
        self.assertTrue(any("stale" in e for e in errors))

    def test_validate_manifest_missing_row(self):
        manifest = empty_manifest("slug")
        manifest["rows"] = []
        url = "https://www.perplexity.ai/search/bacb0d15-2888-45ab-ab32-b30aae743709"
        result = validate_manifest_for_urls(manifest, [url])
        self.assertFalse(result.ok)
        self.assertTrue(any("no manifest row" in e for e in result.errors))

    def test_validate_manifest_url_mismatch(self):
        manifest = empty_manifest("slug")
        manifest["rows"] = [self._row()]
        other = "https://www.perplexity.ai/search/17ab620e-87db-44a6-9f71-c2df9fe63058"
        result = validate_manifest_for_urls(manifest, [other])
        self.assertFalse(result.ok)

    def test_validate_manifest_success(self):
        manifest = empty_manifest("slug")
        manifest["rows"] = [self._row()]
        url = "https://www.perplexity.ai/search/bacb0d15-2888-45ab-ab32-b30aae743709"
        result = validate_manifest_for_urls(manifest, [url])
        self.assertTrue(result.ok)
        self.assertEqual(result.errors, [])


    def test_urls_from_proposal(self):
        from explore_share import urls_from_proposal

        text = """## reviews
reviews:
  - candidate_id: a
    url: https://www.perplexity.ai/search/11111111-1111-4111-8111-111111111111
    verdict: approve
"""
        urls = urls_from_proposal(text)
        self.assertEqual(len(urls), 1)
        self.assertIn("11111111", urls[0])


if __name__ == "__main__":
    unittest.main()
