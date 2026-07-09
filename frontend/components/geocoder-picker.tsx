"use client";

import { useState, useCallback, useEffect } from "react";
import { YMaps, Map, Placemark } from "@pbe/react-yandex-maps";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const API_KEY = "f";

async function geocode(address: string): Promise<{ lat: number; lng: number; display: string } | null> {
  const url = `https://geocode-maps.yandex.ru/1.x/?apikey=${API_KEY}&geocode=${encodeURIComponent(address)}&format=json&results=1&lang=ru_RU`;
  const res = await fetch(url);
  const data = await res.json();
  const members = data?.response?.GeoObjectCollection?.featureMember;
  if (!members || members.length === 0) return null;
  const pos = members[0].GeoObject.Point.pos;
  const [lng, lat] = pos.split(" ").map(Number);
  const display = members[0].GeoObject.metaDataProperty?.GeocoderMetaData?.text || address;
  return { lat, lng, display };
}

interface GeocoderPickerProps {
  value: string;
  onChange: (coords: string) => void;
  city?: string;
  street?: string;
  house?: string;
}

export default function GeocoderPicker({ value, onChange, city, street, house }: GeocoderPickerProps) {
  const parseCoords = (v: string): [number, number] | null => {
    const parts = v.replace(/\s/g, "").split(",").map(Number);
    if (parts.length === 2 && !isNaN(parts[0]) && !isNaN(parts[1])) return [parts[0], parts[1]];
    return null;
  };

  const [marker, setMarker] = useState<[number, number] | null>(() => parseCoords(value));
  const [center, setCenter] = useState<[number, number]>(marker || [55.75254, 37.623082]);
  const [zoom, setZoom] = useState(marker ? 14 : 10);
  const [mapRef, setMapRef] = useState<any>(null);
  const [searchError, setSearchError] = useState("");

  useEffect(() => {
    const coords = parseCoords(value);
    if (coords) {
      setMarker(coords);
      setCenter(coords);
      setZoom(14);
      if (mapRef) {
        setTimeout(() => {
          mapRef.setCenter(coords, 14);
          if (mapRef.container) {
            mapRef.container.fitToViewport();
          }
        }, 100);
      }
    }
  }, [value, mapRef]);

  useEffect(() => {
    if (mapRef) {
      setTimeout(() => {
        if (mapRef.container) {
          mapRef.container.fitToViewport();
        }
      }, 300);
    }
  }, [mapRef]);

  const handleMapClick = useCallback((e: any) => {
    const coords = e.get("coords");
    setMarker(coords);
    setCenter(coords);
    onChange(`${coords[0]},${coords[1]}`);
    setSearchError("");
  }, [onChange]);

  const handleSearch = async () => {
    setSearchError("");
    const parts = [city, street, house].filter(Boolean);
    if (parts.length === 0) {
      setSearchError("Укажите город, улицу и дом для поиска");
      return;
    }
    const address = parts.join(", ");
    const result = await geocode(address);
    if (!result) {
      setSearchError("Адрес не найден. Проверьте данные или укажите координаты вручную на карте.");
      return;
    }
    const coords: [number, number] = [result.lat, result.lng];
    setMarker(coords);
    setCenter(coords);
    setZoom(14);
    onChange(`${result.lat}, ${result.lng}`);
    if (mapRef) {
      setTimeout(() => {
        mapRef.setCenter(coords, 14);
        if (mapRef.container) {
          mapRef.container.fitToViewport();
        }
      }, 100);
    }
  };

  const handleCoordsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    onChange(val);
    const parts = val.replace(/\s/g, "").split(",").map(Number);
    if (parts.length === 2 && !isNaN(parts[0]) && !isNaN(parts[1])) {
      setMarker([parts[0], parts[1]]);
      setCenter([parts[0], parts[1]]);
    }
  };

  return (
    <div className="space-y-3 col-span-full">
      <div className="flex items-end gap-2">
        <div className="space-y-2 flex-1">
          <Label htmlFor="coords">Координаты *</Label>
          <Input
            id="coords"
            placeholder="55.7558, 37.6173"
            value={value}
            onChange={handleCoordsChange}
            required
          />
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="h-9 text-xs"
          onClick={handleSearch}
        >
          Найти по адресу
        </Button>
      </div>

      {searchError && (
        <p className="text-xs text-amber-600 font-medium">{searchError}</p>
      )}

      <p className="text-[10px] text-muted-foreground">
        Нажмите на карту чтобы выбрать координаты, или найдите по адресу через кнопку «Найти по адресу»
      </p>

      <div className="w-full h-64 rounded-xl overflow-hidden border border-slate-200">
        <YMaps query={{ lang: "ru_RU", apikey: API_KEY }}>
        <Map
          defaultState={{ center, zoom }}
          instanceRef={(ref: any) => { if (ref && !mapRef) setMapRef(ref); }}
          width="100%"
          height="100%"
          onClick={handleMapClick}
          options={{ suppressMapOpenBlock: true }}
        >
          {marker && (
            <Placemark
              geometry={marker}
              options={{ preset: "islands#violetDotIconWithCaption", draggable: true }}
              instanceRef={(ref: any) => {
                if (ref) {
                  ref.events.add("dragend", (e: any) => {
                    const coords = e.get("target").geometry.getCoordinates();
                    setMarker(coords);
                    onChange(`${coords[0]},${coords[1]}`);
                  });
                }
              }}
            />
          )}
        </Map>
        </YMaps>
      </div>
    </div>
  );
}
