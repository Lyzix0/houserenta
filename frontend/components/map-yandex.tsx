'use client';

import * as React from 'react';
import { YMaps, Map, Placemark } from '@pbe/react-yandex-maps';

interface Property {
  id: string;
  name: string;
  coordinates: string;
  city: string;
  street: string;
  house: string;
  apartment: string;
  balance?: number;
}

interface MapYandexProps {
  properties: Property[];
  className?: string;
  onMarkerClick?: (property: Property) => void;
  selectedPropertyId?: string | null;
}

const parseCoordinates = (coordStr: string): [number, number] | null => {
  if (!coordStr) return null;
  const parts = coordStr.split(',').map(p => parseFloat(p.trim()));
  if (parts.length === 2 && !isNaN(parts[0]) && !isNaN(parts[1])) {
    return [parts[0], parts[1]];
  }
  return null;
};

export default function MapYandex({
  properties,
  className = "w-full h-80 rounded-2xl overflow-hidden border border-slate-200 shadow-xs",
  onMarkerClick,
  selectedPropertyId
}: MapYandexProps) {
  const validProperties = properties
    .map(p => ({ prop: p, coords: parseCoordinates(p.coordinates) }))
    .filter((item): item is { prop: Property; coords: [number, number] } => item.coords !== null);

  let center: [number, number] = [55.75254, 37.623082];
  let zoom = 11;

  if (validProperties.length > 0) {
    if (selectedPropertyId) {
      const selected = validProperties.find(item => item.prop.id === selectedPropertyId);
      if (selected) {
        center = selected.coords;
        zoom = 14;
      } else {
        const sumLat = validProperties.reduce((sum, item) => sum + item.coords[0], 0);
        const sumLng = validProperties.reduce((sum, item) => sum + item.coords[1], 0);
        center = [sumLat / validProperties.length, sumLng / validProperties.length];
      }
    } else {
      const sumLat = validProperties.reduce((sum, item) => sum + item.coords[0], 0);
      const sumLng = validProperties.reduce((sum, item) => sum + item.coords[1], 0);
      center = [sumLat / validProperties.length, sumLng / validProperties.length];
      if (validProperties.length === 1) {
        zoom = 13;
      }
    }
  }

  return (
    <div className={className}>
      <YMaps query={{ lang: 'ru_RU', apikey: '' }}>
        <Map 
          state={{ center, zoom }} 
          width="100%" 
          height="100%"
          options={{
            suppressMapOpenBlock: true,
            
          }}
        >
          {validProperties.map(({ prop, coords }) => {
            const isSelected = selectedPropertyId === prop.id;
            return (
              <Placemark
                key={prop.id}
                geometry={coords}
                properties={{
                  hintContent: prop.name,
                  balloonContentHeader: `<div style="font-weight: bold; font-family: sans-serif; color: #1e293b;">${prop.name}</div>`,
                  balloonContentBody: `
                    <div style="font-family: sans-serif; font-size: 12px; color: #475569; margin-top: 4px;">
                      ${prop.city}, ул. ${prop.street}, д. ${prop.house}${prop.apartment ? `, кв. ${prop.apartment}` : ''}
                    </div>
                  `,
                  iconCaption: prop.name
                }}
                options={{
                  preset: isSelected 
                    ? 'islands#violetDotIconWithCaption' 
                    : 'islands#tealDotIconWithCaption'
                }}
                onClick={() => {
                  if (onMarkerClick) {
                    onMarkerClick(prop);
                  }
                }}
              />
            );
          })}
        </Map>
      </YMaps>
    </div>
  );
}
