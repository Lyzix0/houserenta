"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface LeaseFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  propertyName: string;
  tenantName: string;
  tenantDoc: string;
  tenantPhone: string;
  monthsOfRent: string; setMonthsOfRent: (v: string) => void;
  rentPrice: string; setRentPrice: (v: string) => void;
  paymentDay: string; setPaymentDay: (v: string) => void;
  readingDay: string; setReadingDay: (v: string) => void;
  tenantError: string;
  unlinkedTenants?: any[];
  selectedTenantUserId?: string;
  setSelectedTenantUserId?: (v: string) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export default function LeaseFormDialog({
  open,
  onOpenChange,
  propertyName,
  tenantName,
  tenantDoc,
  tenantPhone,
  monthsOfRent, setMonthsOfRent,
  rentPrice, setRentPrice,
  paymentDay, setPaymentDay,
  readingDay, setReadingDay,
  tenantError,
  unlinkedTenants = [],
  selectedTenantUserId = "",
  setSelectedTenantUserId = () => {},
  onSubmit,
}: LeaseFormDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-full max-w-md max-h-[90dvh] flex flex-col" showCloseButton>
        <DialogHeader className="shrink-0">
          <DialogTitle>Оформить заселение жильца</DialogTitle>
          <DialogDescription>Свяжите зарегистрированного арендатора с квартирой: {propertyName}</DialogDescription>
        </DialogHeader>
        <div className="overflow-y-auto">
          <form onSubmit={onSubmit} className="space-y-4">
            {tenantName ? (
              <div className="bg-muted/30 p-3 rounded-2xl space-y-1.5 text-xs text-muted-foreground border border-border/50">
                <p><strong>ФИО жильца:</strong> {tenantName}</p>
                <p><strong>Паспорт:</strong> {tenantDoc}</p>
                <p><strong>Телефон:</strong> {tenantPhone}</p>
              </div>
            ) : (
              unlinkedTenants && unlinkedTenants.length > 0 && (
                <div className="space-y-2">
                  <Label htmlFor="tenantSelect">Выберите зарегистрированного жильца *</Label>
                  <Select value={selectedTenantUserId} onValueChange={setSelectedTenantUserId}>
                    <SelectTrigger id="tenantSelect" className="w-full bg-card">
                      <SelectValue placeholder="Выберите жильца..." />
                    </SelectTrigger>
                    <SelectContent>
                      {unlinkedTenants.map((t: any) => (
                        <SelectItem key={t.id} value={t.id}>
                          {t.name} ({t.phone || t.email})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="tMonths">Срок аренды (мес) *</Label>
                <Input id="tMonths" type="number" placeholder="Напр: 12" value={monthsOfRent} onChange={e => setMonthsOfRent(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="tPrice">Месячная рента, ₽ *</Label>
                <Input id="tPrice" type="number" placeholder="35000" value={rentPrice} onChange={e => setRentPrice(e.target.value)} required />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="payDay">День оплаты (1-28) *</Label>
                <Select value={paymentDay} onValueChange={setPaymentDay}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="День оплаты" />
                  </SelectTrigger>
                  <SelectContent>
                    {Array.from({ length: 28 }, (_, i) => String(i + 1)).map(d => (
                      <SelectItem key={d} value={d}>{d}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="readDay">День показаний (1-28) *</Label>
                <Select value={readingDay} onValueChange={setReadingDay}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="День показаний" />
                  </SelectTrigger>
                  <SelectContent>
                    {Array.from({ length: 28 }, (_, i) => String(i + 1)).map(d => (
                      <SelectItem key={d} value={d}>{d}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {tenantError && <p className="text-sm text-red-500 font-semibold">{tenantError}</p>}

            <div className="flex justify-end gap-2 pt-2 border-t">
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>Отмена</Button>
              <Button type="submit">Связать и заселить</Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}
