"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import GeocoderPicker from "@/components/geocoder-picker";

interface PropertyFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editingProperty: any;
  propName: string; setPropName: (v: string) => void;
  coords: string; setCoords: (v: string) => void;
  country: string; setCountry: (v: string) => void;
  region: string; setRegion: (v: string) => void;
  city: string; setCity: (v: string) => void;
  street: string; setStreet: (v: string) => void;
  house: string; setHouse: (v: string) => void;
  apartment: string; setApartment: (v: string) => void;
  gvsTariff: string; setGvsTariff: (v: string) => void;
  hvsTariff: string; setHvsTariff: (v: string) => void;
  el1Tariff: string; setEl1Tariff: (v: string) => void;
  el2Tariff: string; setEl2Tariff: (v: string) => void;
  formError: string;
  onSubmit: (e: React.FormEvent) => void;
}

export default function PropertyFormDialog({
  open,
  onOpenChange,
  editingProperty,
  propName, setPropName,
  coords, setCoords,
  country, setCountry,
  region, setRegion,
  city, setCity,
  street, setStreet,
  house, setHouse,
  apartment, setApartment,
  gvsTariff, setGvsTariff,
  hvsTariff, setHvsTariff,
  el1Tariff, setEl1Tariff,
  el2Tariff, setEl2Tariff,
  formError,
  onSubmit,
}: PropertyFormDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-full max-w-lg max-h-[90dvh] flex flex-col" showCloseButton>
        <DialogHeader className="shrink-0">
          <DialogTitle>{editingProperty ? "Редактировать объект" : "Добавить новый объект"}</DialogTitle>
          <DialogDescription>Заполните параметры арендуемой недвижимости</DialogDescription>
        </DialogHeader>
        <div className="overflow-y-auto">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2 col-span-full">
                <Label htmlFor="pName">Имя объекта (необязательно)</Label>
                <Input id="pName" placeholder="Напр: Студия на Чистых Прудах" value={propName} onChange={e => setPropName(e.target.value)} />
                <p className="text-[10px] text-muted-foreground">Если оставить пустым, имя составится по шаблону: &quot;Улица + число&quot;</p>
              </div>

              <GeocoderPicker value={coords} onChange={setCoords} city={city} street={street} house={house} />

              <div className="space-y-2">
                <Label htmlFor="country">Страна</Label>
                <Input id="country" placeholder="Россия" value={country} onChange={e => setCountry(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="region">Регион *</Label>
                <Input id="region" placeholder="Московская область" value={region} onChange={e => setRegion(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="city">Город *</Label>
                <Input id="city" placeholder="Москва" value={city} onChange={e => setCity(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="street">Улица *</Label>
                <Input id="street" placeholder="Ленина" value={street} onChange={e => setStreet(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="house">Дом *</Label>
                <Input id="house" placeholder="24" value={house} onChange={e => setHouse(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="apartment">Квартира *</Label>
                <Input id="apartment" placeholder="45" value={apartment} onChange={e => setApartment(e.target.value)} required />
              </div>
            </div>

            <div className="border-t pt-4">
              <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-3">Тарифы ЖКХ</h4>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
                <div className="space-y-1">
                  <Label htmlFor="tGvs" className="text-[10px]">ГВС, ₽/м³</Label>
                  <Input id="tGvs" type="number" step="0.01" value={gvsTariff} onChange={e => setGvsTariff(e.target.value)} required />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="tHvs" className="text-[10px]">ХВС, ₽/м³</Label>
                  <Input id="tHvs" type="number" step="0.01" value={hvsTariff} onChange={e => setHvsTariff(e.target.value)} required />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="tEl1" className="text-[10px]">ЭЛ1 (Пик)</Label>
                  <Input id="tEl1" type="number" step="0.01" value={el1Tariff} onChange={e => setEl1Tariff(e.target.value)} required />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="tEl2" className="text-[10px]">ЭЛ2 (Ночь)</Label>
                  <Input id="tEl2" type="number" step="0.01" placeholder="Необяз." value={el2Tariff} onChange={e => setEl2Tariff(e.target.value)} />
                </div>
              </div>
            </div>

            {formError && <p className="text-sm text-red-500 font-semibold">{formError}</p>}

            <div className="flex justify-end gap-2 pt-2 border-t">
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>Отмена</Button>
              <Button type="submit">
                {editingProperty ? "Сохранить изменения" : "Добавить объект"}
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}
