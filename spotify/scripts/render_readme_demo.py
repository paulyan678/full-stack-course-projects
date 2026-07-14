#!/usr/bin/env python3
"""Render the README walkthrough from the Android UI states and local fixtures.

This is intentionally a source-driven walkthrough, not an emulator recording.
The renderer reads the same ``feed.json`` and ``playlists.json`` used by the
Ktor backend, mirrors the Compose screen hierarchy, and reproduces the cover
palette defined by ``CoverFixture.kt``.

Pillow is the only non-stdlib dependency::

    python3 spotify/scripts/render_readme_demo.py
"""

from __future__ import annotations

import argparse
import json
import math
from pathlib import Path
from typing import Sequence

from PIL import Image, ImageDraw, ImageFilter, ImageFont


WIDTH = 960
HEIGHT = 640

BACKGROUND = "#090909"
SURFACE = "#181818"
RAISED = "#282828"
GREEN = "#1ED760"
WHITE = "#FFFFFF"
MUTED = "#B3B3B3"
SUBTLE = "#777777"

# Matches spotify/backend/.../CoverFixture.kt.
COVER_PALETTES = (
    ("#5B21B6", "#EC4899"),
    ("#064E3B", "#34D399"),
    ("#7C2D12", "#FB923C"),
    ("#1E3A8A", "#38BDF8"),
    ("#701A75", "#E879F9"),
)

REPO_ROOT = Path(__file__).resolve().parents[2]
SPOTIFY_ROOT = REPO_ROOT / "spotify"
FEED_PATH = SPOTIFY_ROOT / "backend/src/main/resources/feed.json"
PLAYLISTS_PATH = SPOTIFY_ROOT / "backend/src/main/resources/playlists.json"
DEFAULT_OUTPUT_DIR = REPO_ROOT / "docs/assets/demos"


def load_json(path: Path):
    with path.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def load_fixture_story() -> tuple[list[dict], dict, list[dict]]:
    sections = load_json(FEED_PATH)
    playlists = load_json(PLAYLISTS_PATH)
    albums = [album for section in sections for album in section["albums"]]
    midnight_drive = next(album for album in albums if album["album"] == "Midnight Drive")
    playlist = next(item for item in playlists if item["id"] == midnight_drive["id"])
    night_signals = next(song for song in playlist["songs"] if song["name"] == "Night Signals")

    # Fail loudly if the walkthrough ever drifts from the product fixtures.
    assert midnight_drive["cover"] == f"/covers/{midnight_drive['id']}.svg"
    assert night_signals["src"] == "/songs/night-signals.wav"
    assert night_signals["length"] == "00:05"
    return sections, midnight_drive, playlist["songs"]


def find_font(size: int, *, bold: bool = False) -> ImageFont.FreeTypeFont | ImageFont.ImageFont:
    candidates = (
        "/System/Library/Fonts/SFNS.ttf",
        "/System/Library/Fonts/SFCompact.ttf",
        "/System/Library/Fonts/Supplemental/Arial Bold.ttf" if bold else "/System/Library/Fonts/Supplemental/Arial.ttf",
        "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf" if bold else "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
        "/usr/share/fonts/truetype/liberation2/LiberationSans-Bold.ttf" if bold else "/usr/share/fonts/truetype/liberation2/LiberationSans-Regular.ttf",
    )
    for path in candidates:
        try:
            return ImageFont.truetype(path, size=size)
        except OSError:
            continue
    return ImageFont.load_default()


FONTS = {
    "eyebrow": find_font(13, bold=True),
    "step": find_font(13, bold=True),
    "display": find_font(42, bold=True),
    "body": find_font(18),
    "body_bold": find_font(18, bold=True),
    "small": find_font(13),
    "small_bold": find_font(13, bold=True),
    "tiny": find_font(11),
    "phone_title": find_font(22, bold=True),
    "phone_section": find_font(16, bold=True),
    "phone_body": find_font(12),
    "phone_body_bold": find_font(12, bold=True),
    "phone_small": find_font(10),
    "cover_label": find_font(12, bold=True),
    "cover_number": find_font(28, bold=True),
}


def rgb(hex_color: str) -> tuple[int, int, int]:
    value = hex_color.lstrip("#")
    return tuple(int(value[index : index + 2], 16) for index in (0, 2, 4))


def mix(first: str | tuple[int, int, int], second: str | tuple[int, int, int], amount: float):
    first_rgb = rgb(first) if isinstance(first, str) else first
    second_rgb = rgb(second) if isinstance(second, str) else second
    return tuple(round(a + (b - a) * amount) for a, b in zip(first_rgb, second_rgb))


def rounded_mask(size: tuple[int, int], radius: int) -> Image.Image:
    mask = Image.new("L", size, 0)
    ImageDraw.Draw(mask).rounded_rectangle((0, 0, size[0] - 1, size[1] - 1), radius=radius, fill=255)
    return mask


def draw_gradient(
    size: tuple[int, int],
    start: str,
    end: str,
    *,
    horizontal_bias: float = 0.55,
) -> Image.Image:
    width, height = size
    start_rgb = rgb(start)
    end_rgb = rgb(end)
    image = Image.new("RGB", size)
    pixels = image.load()
    denominator = max(1.0, width * horizontal_bias + height * (1.0 - horizontal_bias))
    for y in range(height):
        for x in range(width):
            t = (x * horizontal_bias + y * (1.0 - horizontal_bias)) / denominator
            pixels[x, y] = mix(start_rgb, end_rgb, min(1.0, t))
    return image


def cover_image(album_id: int, size: int, *, radius: int = 9) -> Image.Image:
    start, end = COVER_PALETTES[(album_id - 1) % len(COVER_PALETTES)]
    cover = draw_gradient((size, size), start, end)
    overlay = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(overlay)
    center = size // 2
    ring_radius = round(size * 0.30)
    ring_width = max(5, round(size * 0.057))
    draw.ellipse(
        (center - ring_radius, center - ring_radius, center + ring_radius, center + ring_radius),
        outline=(255, 255, 255, 46),
        width=ring_width,
    )
    dot_radius = max(4, round(size * 0.09))
    draw.ellipse(
        (center - dot_radius, center - dot_radius, center + dot_radius, center + dot_radius),
        fill=(255, 255, 255, 224),
    )
    label_font = FONTS["cover_label"] if size >= 100 else FONTS["tiny"]
    draw.text((max(8, size * 0.07), max(8, size * 0.07)), "LOCAL SESSIONS", font=label_font, fill=WHITE)
    number = f"{album_id:02d}"
    number_font = FONTS["cover_number"] if size >= 100 else FONTS["small_bold"]
    number_bbox = draw.textbbox((0, 0), number, font=number_font)
    draw.text(
        (max(8, size * 0.07), size - (number_bbox[3] - number_bbox[1]) - max(9, size * 0.08)),
        number,
        font=number_font,
        fill=WHITE,
    )
    composed = Image.alpha_composite(cover.convert("RGBA"), overlay)
    composed.putalpha(rounded_mask((size, size), radius))
    return composed


def text_width(draw: ImageDraw.ImageDraw, text: str, font: ImageFont.ImageFont) -> int:
    bbox = draw.textbbox((0, 0), text, font=font)
    return bbox[2] - bbox[0]


def wrap_text(draw: ImageDraw.ImageDraw, text: str, font: ImageFont.ImageFont, width: int) -> list[str]:
    words = text.split()
    lines: list[str] = []
    current = ""
    for word in words:
        candidate = f"{current} {word}".strip()
        if current and text_width(draw, candidate, font) > width:
            lines.append(current)
            current = word
        else:
            current = candidate
    if current:
        lines.append(current)
    return lines


def draw_wrapped_text(
    draw: ImageDraw.ImageDraw,
    xy: tuple[int, int],
    text: str,
    font: ImageFont.ImageFont,
    fill,
    width: int,
    *,
    spacing: int = 6,
    max_lines: int | None = None,
) -> int:
    lines = wrap_text(draw, text, font, width)
    if max_lines is not None:
        lines = lines[:max_lines]
    x, y = xy
    bbox = draw.textbbox((0, 0), "Ag", font=font)
    line_height = bbox[3] - bbox[1]
    for line in lines:
        draw.text((x, y), line, font=font, fill=fill)
        y += line_height + spacing
    return y


def draw_heart(draw: ImageDraw.ImageDraw, center: tuple[int, int], size: int, *, filled: bool) -> None:
    cx, cy = center
    points: list[tuple[float, float]] = []
    for index in range(120):
        t = 2 * math.pi * index / 120
        x = 16 * math.sin(t) ** 3
        y = 13 * math.cos(t) - 5 * math.cos(2 * t) - 2 * math.cos(3 * t) - math.cos(4 * t)
        points.append((cx + x * size / 36, cy - y * size / 36))
    if filled:
        draw.polygon(points, fill=GREEN)
    else:
        draw.line(points + [points[0]], fill=WHITE, width=max(2, size // 10), joint="curve")


def draw_play(draw: ImageDraw.ImageDraw, center: tuple[int, int], size: int, *, pause: bool = False, color=WHITE) -> None:
    cx, cy = center
    if pause:
        bar = max(3, size // 5)
        gap = max(3, size // 7)
        draw.rounded_rectangle((cx - gap - bar, cy - size // 2, cx - gap, cy + size // 2), radius=2, fill=color)
        draw.rounded_rectangle((cx + gap, cy - size // 2, cx + gap + bar, cy + size // 2), radius=2, fill=color)
    else:
        draw.polygon(
            ((cx - size // 3, cy - size // 2), (cx - size // 3, cy + size // 2), (cx + size // 2, cy)),
            fill=color,
        )


def draw_home_icon(draw: ImageDraw.ImageDraw, center: tuple[int, int], *, active: bool) -> None:
    cx, cy = center
    color = WHITE if active else MUTED
    draw.polygon(((cx - 9, cy), (cx, cy - 8), (cx + 9, cy), (cx + 7, cy), (cx + 7, cy + 8), (cx - 7, cy + 8), (cx - 7, cy)), fill=color)
    if not active:
        draw.rectangle((cx - 4, cy + 2, cx + 4, cy + 8), fill=SURFACE)


def draw_nav_bar(draw: ImageDraw.ImageDraw, box: tuple[int, int, int, int], *, active: str) -> None:
    x0, y0, x1, y1 = box
    draw.rectangle(box, fill=SURFACE)
    home_x = x0 + (x1 - x0) // 3
    favorite_x = x0 + 2 * (x1 - x0) // 3
    icon_y = y0 + 15
    draw_home_icon(draw, (home_x, icon_y), active=active == "home")
    draw_heart(draw, (favorite_x, icon_y), 18, filled=active == "favorites")
    draw.text((home_x - text_width(draw, "Home", FONTS["phone_small"]) // 2, icon_y + 13), "Home", font=FONTS["phone_small"], fill=WHITE if active == "home" else MUTED)
    favorite_label = "Favorites"
    draw.text((favorite_x - text_width(draw, favorite_label, FONTS["phone_small"]) // 2, icon_y + 13), favorite_label, font=FONTS["phone_small"], fill=WHITE if active == "favorites" else MUTED)


def draw_player_bar(
    screen: Image.Image,
    box: tuple[int, int, int, int],
    album: dict,
    song: dict,
    *,
    progress: float,
) -> None:
    draw = ImageDraw.Draw(screen)
    x0, y0, x1, y1 = box
    draw.rounded_rectangle(box, radius=10, fill=RAISED)
    cover_size = 48
    screen.alpha_composite(cover_image(album["id"], cover_size, radius=6), (x0 + 8, y0 + 8))
    draw.text((x0 + 66, y0 + 10), song["name"], font=FONTS["phone_body_bold"], fill=WHITE)
    draw.text((x0 + 66, y0 + 29), album["artists"], font=FONTS["phone_small"], fill=MUTED)
    draw_play(draw, (x1 - 27, y0 + 31), 18, pause=True)
    track_left = x0 + 10
    track_right = x1 - 10
    track_y = y1 - 7
    draw.rounded_rectangle((track_left, track_y, track_right, track_y + 3), radius=2, fill="#555555")
    active_x = round(track_left + (track_right - track_left) * progress)
    draw.rounded_rectangle((track_left, track_y, active_x, track_y + 3), radius=2, fill=GREEN)
    draw.ellipse((active_x - 3, track_y - 2, active_x + 3, track_y + 5), fill=GREEN)


def phone_shell() -> tuple[Image.Image, tuple[int, int, int, int]]:
    phone_w, phone_h = 356, 612
    phone = Image.new("RGBA", (phone_w, phone_h), (0, 0, 0, 0))
    shadow = Image.new("RGBA", (phone_w + 40, phone_h + 40), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    shadow_draw.rounded_rectangle((20, 16, phone_w + 19, phone_h + 15), radius=43, fill=(0, 0, 0, 205))
    shadow = shadow.filter(ImageFilter.GaussianBlur(16))
    phone.alpha_composite(shadow.crop((20, 20, phone_w + 20, phone_h + 20)), (0, 0))
    draw = ImageDraw.Draw(phone)
    draw.rounded_rectangle((1, 1, phone_w - 2, phone_h - 2), radius=40, fill="#101010", outline="#424242", width=2)
    screen_box = (10, 10, phone_w - 11, phone_h - 11)
    draw.rounded_rectangle(screen_box, radius=32, fill=BACKGROUND)
    draw.rounded_rectangle((phone_w // 2 - 38, 17, phone_w // 2 + 38, 23), radius=4, fill="#303030")
    return phone, screen_box


def render_home_phone(sections: Sequence[dict]) -> Image.Image:
    phone, screen_box = phone_shell()
    screen = Image.new("RGBA", (screen_box[2] - screen_box[0] + 1, screen_box[3] - screen_box[1] + 1), BACKGROUND)
    draw = ImageDraw.Draw(screen)
    width, height = screen.size
    draw.text((16, 8), "9:41", font=FONTS["phone_small"], fill=WHITE)
    draw.text((18, 41), "Good listening", font=FONTS["phone_title"], fill=WHITE)
    draw_wrapped_text(draw, (18, 73), "No account required — every track comes from your local server.", FONTS["phone_body"], MUTED, width - 36, spacing=2, max_lines=2)

    section = sections[0]
    draw.text((18, 116), section["section_title"], font=FONTS["phone_section"], fill=WHITE)
    card_y = 143
    card_size = 132
    gap = 12
    for index, album in enumerate(section["albums"][:2]):
        x = 18 + index * (card_size + gap)
        screen.alpha_composite(cover_image(album["id"], card_size), (x, card_y))
        draw.text((x, card_y + card_size + 7), album["album"], font=FONTS["phone_body_bold"], fill=WHITE)
        draw.text((x, card_y + card_size + 25), album["artists"], font=FONTS["phone_small"], fill=MUTED)

    second = sections[1]
    draw.text((18, 323), second["section_title"], font=FONTS["phone_section"], fill=WHITE)
    small_size = 84
    for index, album in enumerate(second["albums"][:3]):
        x = 18 + index * (small_size + 13)
        screen.alpha_composite(cover_image(album["id"], small_size, radius=7), (x, 350))
        label = album["album"].split()[0]
        draw.text((x, 440), label, font=FONTS["phone_small"], fill=WHITE)

    draw_nav_bar(draw, (0, height - 50, width, height), active="home")
    screen.putalpha(rounded_mask(screen.size, 31))
    phone.alpha_composite(screen, (screen_box[0], screen_box[1]))
    return phone


def render_playlist_phone(album: dict, songs: Sequence[dict], *, playing: bool) -> Image.Image:
    phone, screen_box = phone_shell()
    screen = Image.new("RGBA", (screen_box[2] - screen_box[0] + 1, screen_box[3] - screen_box[1] + 1), BACKGROUND)
    draw = ImageDraw.Draw(screen)
    width, height = screen.size
    draw.text((16, 8), "9:41", font=FONTS["phone_small"], fill=WHITE)

    cover_height = 206 if not playing else 158
    cover = cover_image(album["id"], 335, radius=0).crop((0, 32, 335, 32 + cover_height))
    screen.alpha_composite(cover, (0, 27))
    draw.ellipse((12, 39, 42, 69), fill=(0, 0, 0, 130))
    draw.line((32, 47, 22, 54, 32, 62), fill=WHITE, width=3, joint="curve")

    info_y = cover_height + 40
    draw.text((18, info_y), album["album"], font=FONTS["phone_title"], fill=WHITE)
    info_y += 31
    description = album["description"]
    info_y = draw_wrapped_text(draw, (18, info_y), description, FONTS["phone_small"], MUTED, width - 74, spacing=2, max_lines=2)
    draw.text((18, info_y + 4), album["artists"], font=FONTS["phone_body_bold"], fill=WHITE)
    draw.text((18, info_y + 22), album["year"], font=FONTS["phone_small"], fill=MUTED)
    heart_center = (width - 35, info_y + 17)
    draw_heart(draw, heart_center, 28, filled=playing)

    row_y = info_y + 54
    for index, song in enumerate(songs[:3]):
        selected = playing and index == 0
        if selected:
            draw_play(draw, (27, row_y + 18), 14, color=GREEN)
        else:
            draw.text((24, row_y + 9), str(index + 1), font=FONTS["phone_body"], fill=MUTED)
        draw.text((48, row_y + 4), song["name"], font=FONTS["phone_body_bold"], fill=GREEN if selected else WHITE)
        draw.text((48, row_y + 22), song["lyric"], font=FONTS["phone_small"], fill=MUTED)
        draw.text((width - 47, row_y + 13), song["length"], font=FONTS["phone_small"], fill=MUTED)
        row_y += 45
        if row_y > height - (142 if playing else 55):
            break

    if playing:
        player_top = height - 128
        draw_player_bar(screen, (7, player_top, width - 7, player_top + 77), album, songs[0], progress=0.38)
        draw_nav_bar(draw, (0, height - 50, width, height), active="home")
    else:
        draw_nav_bar(draw, (0, height - 50, width, height), active="home")

    screen.putalpha(rounded_mask(screen.size, 31))
    phone.alpha_composite(screen, (screen_box[0], screen_box[1]))
    return phone


def render_favorites_phone(album: dict, song: dict) -> Image.Image:
    phone, screen_box = phone_shell()
    screen = Image.new("RGBA", (screen_box[2] - screen_box[0] + 1, screen_box[3] - screen_box[1] + 1), BACKGROUND)
    draw = ImageDraw.Draw(screen)
    width, height = screen.size
    draw.text((16, 8), "9:41", font=FONTS["phone_small"], fill=WHITE)
    draw.text((18, 47), "Your favorites", font=FONTS["phone_title"], fill=WHITE)

    cover_size = 78
    screen.alpha_composite(cover_image(album["id"], cover_size, radius=7), (18, 91))
    draw.text((112, 105), album["album"], font=FONTS["phone_section"], fill=WHITE)
    draw.text((112, 134), album["artists"], font=FONTS["phone_body"], fill=MUTED)
    draw_heart(draw, (width - 30, 129), 22, filled=True)
    draw.line((18, 187, width - 18, 187), fill="#242424", width=1)
    draw.text((18, 211), "Saved on this device", font=FONTS["phone_body_bold"], fill=WHITE)
    draw.text((18, 235), "Room keeps favorites available between sessions.", font=FONTS["phone_small"], fill=MUTED)

    player_top = height - 128
    draw_player_bar(screen, (7, player_top, width - 7, player_top + 77), album, song, progress=0.72)
    draw_nav_bar(draw, (0, height - 50, width, height), active="favorites")

    screen.putalpha(rounded_mask(screen.size, 31))
    phone.alpha_composite(screen, (screen_box[0], screen_box[1]))
    return phone


def base_canvas() -> Image.Image:
    image = draw_gradient((WIDTH, HEIGHT), "#080808", "#151515", horizontal_bias=0.72).convert("RGBA")
    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.ellipse((548, 42, 1068, 610), fill=(30, 215, 96, 26))
    glow = glow.filter(ImageFilter.GaussianBlur(70))
    image = Image.alpha_composite(image, glow)
    # Deterministic, sparse texture prevents flat GIF banding without visual noise.
    texture = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    texture_draw = ImageDraw.Draw(texture)
    for y in range(0, HEIGHT, 11):
        for x in range((y * 7) % 17, WIDTH, 17):
            texture_draw.point((x, y), fill=(255, 255, 255, 7))
    return Image.alpha_composite(image, texture)


STEPS = (
    (
        "01",
        "Browse local albums",
        "The Home feed is populated by feed.json and rendered as Compose LazyRows — no account or remote catalog required.",
        "Tap “Midnight Drive”",
    ),
    (
        "02",
        "Open Midnight Drive",
        "The playlist route loads the album and its three five-second WAV fixtures from the local Ktor API.",
        "Save the album, then play “Night Signals”",
    ),
    (
        "03",
        "Save it and press play",
        "The heart persists Midnight Drive with Room while Media3 starts /songs/night-signals.wav and exposes playback state.",
        "Night Signals • playing",
    ),
    (
        "04",
        "Find it in Favorites",
        "Favorites observes the Room-backed flow. The app-level mini-player remains available while navigating between tabs.",
        "Midnight Drive • saved on this device",
    ),
)


def draw_step_progress(draw: ImageDraw.ImageDraw, active_index: int) -> None:
    x = 56
    y = 537
    gap = 83
    for index, (_, _, _, _) in enumerate(STEPS):
        active = index <= active_index
        dot_color = GREEN if active else "#474747"
        draw.ellipse((x - 5, y - 5, x + 5, y + 5), fill=dot_color)
        if index < len(STEPS) - 1:
            line_color = GREEN if index < active_index else "#343434"
            draw.rounded_rectangle((x + 9, y - 1, x + gap - 9, y + 2), radius=2, fill=line_color)
        draw.text((x - 7, y + 15), str(index + 1), font=FONTS["tiny"], fill=WHITE if index == active_index else SUBTLE)
        x += gap


def render_scene(index: int, phone: Image.Image) -> Image.Image:
    canvas = base_canvas()
    draw = ImageDraw.Draw(canvas)
    number, heading, body, action = STEPS[index]

    # Persistent provenance badge: the walkthrough cannot be mistaken for footage.
    badge = (54, 37, 342, 67)
    draw.rounded_rectangle(badge, radius=15, fill="#10351F", outline="#1E6C3A", width=1)
    draw.ellipse((68, 49, 76, 57), fill=GREEN)
    draw.text((86, 45), "SOURCE-DRIVEN COMPOSE WALKTHROUGH", font=FONTS["eyebrow"], fill=WHITE)

    draw.text((56, 106), f"STEP {number} / 04", font=FONTS["step"], fill=GREEN)
    draw_wrapped_text(draw, (54, 140), heading, FONTS["display"], WHITE, 420, spacing=2, max_lines=2)
    draw_wrapped_text(draw, (56, 250), body, FONTS["body"], MUTED, 414, spacing=7, max_lines=5)

    action_y = 405
    draw.rounded_rectangle((54, action_y, 458, action_y + 68), radius=16, fill="#171717", outline="#343434", width=1)
    draw.ellipse((72, action_y + 20, 100, action_y + 48), fill=GREEN)
    if index == 2:
        draw_play(draw, (86, action_y + 34), 13, pause=True, color=BACKGROUND)
    elif index == 3:
        draw_heart(draw, (86, action_y + 34), 15, filled=True)
    else:
        draw.ellipse((82, action_y + 30, 90, action_y + 38), fill=BACKGROUND)
        draw.ellipse((77, action_y + 25, 95, action_y + 43), outline=BACKGROUND, width=2)
    draw_wrapped_text(draw, (114, action_y + 15), action, FONTS["small_bold"], WHITE, 324, spacing=2, max_lines=2)

    draw_step_progress(draw, index)
    draw.text((56, 594), "Rendered from Compose UI states + local JSON fixtures", font=FONTS["tiny"], fill=SUBTLE)

    canvas.alpha_composite(phone, (556, 14))
    return canvas.convert("RGB")


def transition_frames(start: Image.Image, end: Image.Image, count: int = 4) -> list[Image.Image]:
    frames: list[Image.Image] = []
    for step in range(1, count + 1):
        amount = step / (count + 1)
        # Smoothstep keeps both UI states readable at the ends of the transition.
        eased = amount * amount * (3 - 2 * amount)
        frames.append(Image.blend(start, end, eased))
    return frames


def build_animation(scenes: Sequence[Image.Image]) -> tuple[list[Image.Image], list[int]]:
    frames: list[Image.Image] = []
    durations: list[int] = []
    for index, scene in enumerate(scenes):
        frames.append(scene)
        durations.append(1650 if index < len(scenes) - 1 else 2300)
        if index < len(scenes) - 1:
            transitions = transition_frames(scene, scenes[index + 1])
            frames.extend(transitions)
            durations.extend([95] * len(transitions))
    return frames, durations


def save_gif(frames: Sequence[Image.Image], durations: Sequence[int], path: Path) -> None:
    # A single adaptive palette keeps the file compact and prevents palette flash.
    contact = Image.new("RGB", (WIDTH, HEIGHT * len(frames)), BACKGROUND)
    for index, frame in enumerate(frames):
        contact.paste(frame, (0, HEIGHT * index))
    palette_source = contact.quantize(colors=112, method=Image.Quantize.MEDIANCUT)
    palette = palette_source.crop((0, 0, WIDTH, HEIGHT))
    quantized = [
        frame.quantize(palette=palette, dither=Image.Dither.NONE)
        for frame in frames
    ]
    quantized[0].save(
        path,
        save_all=True,
        append_images=quantized[1:],
        duration=list(durations),
        loop=0,
        optimize=True,
        disposal=1,
    )


def render(output_dir: Path) -> tuple[Path, Path, int]:
    sections, album, songs = load_fixture_story()
    phones = (
        render_home_phone(sections),
        render_playlist_phone(album, songs, playing=False),
        render_playlist_phone(album, songs, playing=True),
        render_favorites_phone(album, songs[0]),
    )
    scenes = [render_scene(index, phone) for index, phone in enumerate(phones)]
    frames, durations = build_animation(scenes)

    output_dir.mkdir(parents=True, exist_ok=True)
    gif_path = output_dir / "spotify-demo.gif"
    poster_path = output_dir / "spotify-poster.png"
    scenes[-1].save(poster_path, format="PNG", optimize=True)
    save_gif(frames, durations, gif_path)
    return gif_path, poster_path, len(frames)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--output-dir",
        type=Path,
        default=DEFAULT_OUTPUT_DIR,
        help=f"Asset destination (default: {DEFAULT_OUTPUT_DIR})",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    gif_path, poster_path, frame_count = render(args.output_dir.resolve())
    print(f"Rendered {gif_path} ({WIDTH}x{HEIGHT}, {frame_count} frames)")
    print(f"Rendered {poster_path} ({WIDTH}x{HEIGHT})")


if __name__ == "__main__":
    main()
