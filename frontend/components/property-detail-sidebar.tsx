"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";

interface PropertyDetailSidebarProps {
  property: any;
  user: any;
  onClose: () => void;
  onContractClick: () => void;
  onRemoveTenant: (id: string) => void;
  onSubmitReadings: (e: React.FormEvent) => void;
  onLogPayment: (e: React.FormEvent) => void;
  onAddCustomItem: (e: React.FormEvent) => void;
  onShowTenantModal: (app: any) => void;
  onManualLeaseClick?: () => void;
  onBillClick: (bill: any) => void;
  gvsVal: string;
  setGvsVal: (v: string) => void;
  hvsVal: string;
  setHvsVal: (v: string) => void;
  el1Val: string;
  setEl1Val: (v: string) => void;
  el2Val: string;
  setEl2Val: (v: string) => void;
  paymentAmount: string;
  setPaymentAmount: (v: string) => void;
  customDesc: string;
  setCustomDesc: (v: string) => void;
  customAmount: string;
  setCustomAmount: (v: string) => void;
}

export default function PropertyDetailSidebar({
  property,
  user,
  onClose,
  onContractClick,
  onRemoveTenant,
  onSubmitReadings,
  onLogPayment,
  onAddCustomItem,
  onShowTenantModal,
  onManualLeaseClick,
  onBillClick,
  gvsVal, setGvsVal,
  hvsVal, setHvsVal,
  el1Val, setEl1Val,
  el2Val, setEl2Val,
  paymentAmount, setPaymentAmount,
  customDesc, setCustomDesc,
  customAmount, setCustomAmount,
}: PropertyDetailSidebarProps) {
  const router = useRouter();

  return (
    <div className="fixed inset-0 z-40 flex items-center justify-end bg-black/60 backdrop-blur-xs">
      <div className="bg-background w-full max-w-2xl h-full p-6 shadow-2xl overflow-y-auto flex flex-col justify-between border-l border-border">
        <div>
          <div className="flex justify-between items-start border-b pb-4 mb-6 border-border">
            <div>
              <h2 className="text-xl font-bold">{property.name}</h2>
              <p className="text-xs text-muted-foreground mt-1">
                {property.country || ""}, {property.region}, {property.city}, ул. {property.street}, д. {property.house}, кв. {property.apartment}
              </p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-8 w-8 rounded-full"
            >
              ✕
            </Button>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6 bg-muted/30 border border-border/50 p-4 rounded-2xl text-xs">
            <div>
              <p className="text-muted-foreground font-semibold">Тариф ГВС</p>
              <p className="font-bold mt-0.5">{property.gvs_tariff} ₽/м³</p>
            </div>
            <div>
              <p className="text-muted-foreground font-semibold">Тариф ХВС</p>
              <p className="font-bold mt-0.5">{property.hvs_tariff} ₽/м³</p>
            </div>
            <div>
              <p className="text-muted-foreground font-semibold">Тариф ЭЛ1</p>
              <p className="font-bold mt-0.5">{property.el1_tariff} ₽/кВт*ч</p>
            </div>
            <div>
              <p className="text-muted-foreground font-semibold">Тариф ЭЛ2</p>
              <p className="font-bold mt-0.5">{property.el2_tariff ? `${property.el2_tariff} ₽/кВт*ч` : "Не исп."}</p>
            </div>
          </div>

          <div className="mb-6 flex justify-between items-center bg-card border border-border p-4 rounded-2xl">
            <div>
              <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Текущий баланс</h4>
              <p className="text-xs text-muted-foreground mt-0.5">Взаиморасчеты с арендатором</p>
            </div>
            <p className={`text-2xl font-bold tracking-tight ${property.balance < 0 ? "text-destructive" : "text-emerald-600"}`}>
              {property.balance.toLocaleString("ru-RU")} ₽
            </p>
          </div>

          <div className="border-t pt-6 mb-8 border-border">
            {property.tenant ? (
              <div className="space-y-6">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-base font-bold">Арендатор: {property.tenant.name}</h3>
                    <p className="text-xs text-muted-foreground mt-1">Документ: {property.tenant.document}</p>
                    <p className="text-xs text-muted-foreground">Тел: {property.tenant.phone}</p>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      onClick={onContractClick}
                      size="sm"
                      variant="outline"
                      className="text-xs h-8 font-medium"
                    >
                      Договор
                    </Button>
                    <Button
                      onClick={() => onRemoveTenant(property.id)}
                      size="sm"
                      variant="destructive"
                      className="text-xs h-8 font-medium"
                    >
                      Выселить
                    </Button>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-3">
                    <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Показания</h4>

                    <form onSubmit={onSubmitReadings} className="bg-muted/30 p-3 rounded-2xl space-y-2 border border-border/50">
                      <p className="text-[10px] font-semibold text-muted-foreground">Внести новые показания:</p>
                      <div className="grid grid-cols-3 gap-2">
                        <Input type="number" step="0.01" placeholder="ГВС" className="h-8 text-xs p-1" value={gvsVal} onChange={e => setGvsVal(e.target.value)} required />
                        <Input type="number" step="0.01" placeholder="ХВС" className="h-8 text-xs p-1" value={hvsVal} onChange={e => setHvsVal(e.target.value)} required />
                        <Input type="number" step="0.01" placeholder="Т1 Пик" className="h-8 text-xs p-1" value={el1Val} onChange={e => setEl1Val(e.target.value)} required />
                      </div>
                      {property.el2_tariff !== null && (
                        <Input type="number" step="0.01" placeholder="Т2 Ночь" className="h-8 text-xs p-2" value={el2Val} onChange={e => setEl2Val(e.target.value)} />
                      )}
                      <Button type="submit" size="sm" className="w-full h-8 mt-1">
                        Сохранить показания
                      </Button>
                    </form>

                    <div className="space-y-2 max-h-48 overflow-y-auto">
                      {property.readings.map((reading: any) => (
                        <div key={reading.id} className="p-2 border border-border/50 rounded-xl flex justify-between items-center bg-card text-[11px]">
                          <div>
                            <span className="font-semibold">{new Date(reading.date).toLocaleDateString("ru-RU")}</span>
                            <p className="text-[9px] text-muted-foreground">ГВС:{reading.gvs} ХВС:{reading.hvs} Т1:{reading.el1}</p>
                          </div>
                          <Badge variant={reading.is_accounted ? "secondary" : "outline"} className={reading.is_accounted ? "" : "border-amber-500/50 text-amber-600 bg-amber-500/5"}>
                            {reading.is_accounted ? "Учтено" : "Новые"}
                          </Badge>
                        </div>
                      ))}
                    </div>
                  </div>

                  <div className="space-y-3">
                    <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Счета и платежи</h4>

                    <form onSubmit={onLogPayment} className="bg-muted/30 p-3 rounded-2xl space-y-2 border border-border/50">
                      <p className="text-[10px] font-semibold text-muted-foreground">Зарегистрировать оплату:</p>
                      <div className="flex gap-2">
                        <Input type="number" placeholder="Сумма, руб" className="h-8 text-xs flex-1" value={paymentAmount} onChange={e => setPaymentAmount(e.target.value)} required />
                        <Button type="submit" size="sm" className="h-8 shrink-0 px-2.5">Внести</Button>
                      </div>
                    </form>

                    <form onSubmit={onAddCustomItem} className="bg-muted/30 p-3 rounded-2xl space-y-2 border border-border/50">
                      <p className="text-[10px] font-semibold text-muted-foreground">Пункт в следующий счет:</p>
                      <div className="grid grid-cols-2 gap-2">
                        <Input placeholder="Комментарий" className="h-8 text-xs p-2" value={customDesc} onChange={e => setCustomDesc(e.target.value)} required />
                        <Input type="number" placeholder="Сумма" className="h-8 text-xs p-2" value={customAmount} onChange={e => setCustomAmount(e.target.value)} required />
                      </div>
                      <Button type="submit" size="sm" className="w-full h-8">Добавить в счет</Button>
                    </form>

                    <div className="space-y-2 max-h-48 overflow-y-auto">
                      {property.bills.map((bill: any) => (
                        <div key={bill.id} className="p-2 border border-border/50 rounded-xl bg-card text-[11px] flex justify-between items-center">
                          <div>
                            <p className="font-semibold">{bill.total} ₽ ({bill.type === "rent" ? "Счет" : "Взнос"})</p>
                            <p className="text-[9px] text-muted-foreground">{new Date(bill.date).toLocaleDateString("ru-RU")}</p>
                          </div>
                          <div className="flex items-center gap-1.5">
                            <Badge variant={bill.status === "paid" ? "secondary" : "destructive"}>
                              {bill.status === "paid" ? "Оплачен" : "Ожидает"}
                            </Badge>
                            <Button onClick={() => onBillClick(bill)} size="sm" variant="ghost" className="h-6 px-1.5 text-xs" title="Распечатать">Печать</Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                <div className="text-center py-4 bg-muted/20 rounded-2xl border border-dashed border-border p-6 flex flex-col items-center gap-3">
                  <p className="text-sm font-semibold">Объект свободен</p>
                  <p className="text-xs text-muted-foreground max-w-sm">
                    Вы можете отправить прямую ссылку на заселение будущему жильцу. Он сможет зарегистрироваться и откликнуться автоматически.
                  </p>
                  <div className="flex flex-col sm:flex-row gap-2 w-full max-w-sm justify-center mt-2">
                    <Button
                      type="button"
                      onClick={() => {
                        const link = window.location.origin + "/apply/" + property.id;
                        navigator.clipboard.writeText(link);
                        alert("Ссылка для заселения скопирована в буфер обмена! Отправьте её будущему арендатору.");
                      }}
                      variant="outline"
                      className="text-xs font-semibold rounded-full px-4 h-9"
                    >
                      Поделиться ссылкой
                    </Button>
                    <Button
                      type="button"
                      onClick={() => onShowTenantModal({ tenant_user_id: "", name: "", document: "", phone: "" })}
                      variant="outline"
                      className="text-xs font-semibold rounded-full px-4 h-9"
                    >
                      Заселить вручную
                    </Button>
                  </div>
                </div>

                <div className="border-t border-border pt-6 space-y-4">
                  <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex justify-between items-center">
                    <span>Поступившие отклики ({property.applications?.length || 0})</span>
                  </h4>

                  {(!property.applications || property.applications.length === 0) ? (
                    <p className="text-xs text-muted-foreground italic text-center py-4 bg-muted/20 rounded-2xl border border-dashed border-border">
                      Откликов от жильцов на данный момент нет.
                    </p>
                  ) : (
                    <div className="space-y-3">
                      {property.applications.map((app: any) => (
                        <div key={app.id} className="p-4 border border-border rounded-2xl bg-card flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs">
                          <div>
                            <p className="font-bold text-foreground">{app.name}</p>
                            <p className="text-[10px] text-muted-foreground mt-0.5">Тел: {app.phone}</p>
                            <p className="text-[10px] text-muted-foreground">Документ: {app.document}</p>
                            <p className="text-[9px] text-muted-foreground mt-1">Отклик получен: {new Date(app.date).toLocaleDateString("ru-RU")}</p>
                          </div>
                          <div className="flex gap-2 shrink-0">
                            <Button type="button" onClick={() => router.push("/chat")} variant="outline" className="text-[10px] h-8 px-3 font-semibold">Чат</Button>
                            <Button type="button" onClick={() => onShowTenantModal(app)} className="text-[10px] h-8 px-3">Оформить & Заселить</Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>

        <Button onClick={onClose} variant="secondary" className="w-full mt-6">
          Закрыть
        </Button>
      </div>
    </div>
  );
}
