"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/use-user";
import { useQueryClient } from "@tanstack/react-query";
import { useProperties, useVacantProperties } from "@/hooks/use-properties";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import MapYandex from "@/components/map-yandex";
import { AuthGuard } from "@/components/auth-guard";

import { 
  getProperties, 
  getVacantProperties, 
  payBill, 
  submitReadings, 
  applyForProperty 
} from "@/lib/properties";

const COMPLIMENTS = [
  "Вы отлично управляетесь со своей недвижимостью!",
  "Ваша улыбка освещает этот день ярче солнца!",
  "Прекрасный день для продуктивной работы и отдыха!",
  "Вы делаете этот мир лучше, будучи заботливым арендодателем!",
  "Ваша внимательность к деталям восхищает!",
  "Отличный выбор тарифов ЖКХ, очень экономно!",
  "Ваши арендаторы наверняка считают вас лучшим!",
  "Каждый день — это новая возможность, и вы используете её на все сто!",
];

export default function HomePage() {
  return (
    <AuthGuard mode="auth">
      <DashboardPage />
    </AuthGuard>
  );
}

function DashboardPage() {
  const queryClient = useQueryClient();
  const { data: user, isLoading, refetch } = useUser();

  const [compliment, setCompliment] = useState("");
  const [currentTime, setCurrentTime] = useState("");
  const [isClient, setIsClient] = useState(false);
  const [selectedPropertyId, setSelectedPropertyId] = useState<string | null>(null);

  const [showReadingModal, setShowReadingModal] = useState(false);
  const [gvs, setGvs] = useState("");
  const [hvs, setHvs] = useState("");
  const [el1, setEl1] = useState("");
  const [el2, setEl2] = useState("");

  const [showInvoiceModal, setShowInvoiceModal] = useState(false);
  const [selectedInvoiceBill, setSelectedInvoiceBill] = useState<any | null>(null);

  const [showTopUpModal, setShowTopUpModal] = useState(false);
  const [topUpAmount, setTopUpAmount] = useState("");

  const { data: propertiesData } = useProperties({
    enabled: !!user && (user.role === "landlord" || !!user.tenantPropertyId)
  });

  const { data: vacantPropertiesData } = useVacantProperties({
    enabled: !!user && user.role === "tenant"
  });

  const properties = propertiesData || [];
  const vacantProperties = vacantPropertiesData || [];

  const handleTopUpSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !user.tenantPropertyId || !topUpAmount) return;

    try {
      await payBill(user.tenantPropertyId, { amount: parseFloat(topUpAmount) });
      alert("Баланс успешно пополнен!");
      setShowTopUpModal(false);
      setTopUpAmount("");
      queryClient.invalidateQueries({ queryKey: ["properties"] });
      refetch();
    } catch (err: any) {
      alert(err.message || "Ошибка пополнения");
    }
  };

  const router = useRouter()

  useEffect(() => {
    setIsClient(true);
    setCompliment(COMPLIMENTS[Math.floor(Math.random() * COMPLIMENTS.length)]);
    
    const updateTime = () => {
      const now = new Date();
      setCurrentTime(
        now.toLocaleDateString("ru-RU", {
          day: "numeric",
          month: "long",
          year: "numeric",
        }) + " " + now.toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })
      );
    };
    updateTime();
    const interval = setInterval(updateTime, 60000);
    return () => clearInterval(interval);
  }, []);

  if (isLoading || !isClient) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background text-foreground">
        <div className="flex flex-col items-center gap-2">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
          <p className="text-sm text-muted-foreground">Загрузка личного кабинета...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const handlePayBill = async (propertyId: string, billId: string, total: number) => {
    try {
      await payBill(propertyId, { amount: total, billId });
      alert("Оплата успешно проведена!");
      queryClient.invalidateQueries({ queryKey: ["properties"] });
      refetch();
    } catch (err: any) {
      alert(err.message || "Ошибка оплаты");
    }
  };

  const handleAddReading = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user.tenantPropertyId) return;

    try {
      await submitReadings(user.tenantPropertyId, {
        gvs: parseFloat(gvs),
        hvs: parseFloat(hvs),
        el1: parseFloat(el1),
        el2: el2 ? parseFloat(el2) : null
      });
      alert("Показания счетчиков успешно переданы!");
      setShowReadingModal(false);
      setGvs("");
      setHvs("");
      setEl1("");
      setEl2("");
      queryClient.invalidateQueries({ queryKey: ["properties"] });
    } catch (err: any) {
      alert(err.message || "Ошибка передачи показаний");
    }
  };

  const handleApplyProperty = async (propertyId: string) => {
    try {
      await applyForProperty(propertyId);
      alert("Вы успешно откликнулись на квартиру!");
      queryClient.invalidateQueries({ queryKey: ["vacantProperties"] });
    } catch (err: any) {
      alert(err.message || "Ошибка отправки отклика");
    }
  };

  const isTenant = user.role === "tenant";
  const myProperty = properties.find(p => p.id === user.tenantPropertyId);

  const negBalanceProps = properties.filter(p => p.balance < 0);
  const checkoutProps = properties.filter(p => {
    if (!p.tenant) return false;
    const end = new Date(p.tenant.end_date);
    const diff = end.getTime() - new Date().getTime();
    const diffDays = Math.ceil(diff / (1000 * 60 * 60 * 24));
    return diffDays <= 3;
  });
  const readingsProps = properties.filter(p => {
    if (!p.tenant) return false;
    const pDay = p.tenant.payment_day;
    const paymentDate = new Date();
    paymentDate.setDate(pDay);
    if (paymentDate < new Date()) {
      paymentDate.setMonth(paymentDate.getMonth() + 1);
    }
    const diffDays = Math.ceil((paymentDate.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));
    
    if (diffDays <= 5) {
      const lastRead = p.readings.length > 0 ? p.readings[0] : null;
      if (!lastRead) return true;
      const age = (new Date().getTime() - new Date(lastRead.date).getTime()) / (1000 * 60 * 60 * 24);
      return age > 25;
    }
    return false;
  });
  const recentBills = properties.flatMap((p: any) => 
    p.bills.filter((b: any) => {
      const age = (new Date().getTime() - new Date(b.date).getTime()) / (1000 * 60 * 60 * 24);
      return age <= 30;
    }).map((b: any) => ({ ...b, propertyName: p.name }))
  );

  const totalUrgentCount = negBalanceProps.length + checkoutProps.length + readingsProps.length + recentBills.length;
  const totalProperties = properties.length;
  const occupiedCount = properties.filter((p: any) => p.tenant).length;
  const vacantCount = totalProperties - occupiedCount;
  const totalSaldo = properties.reduce((acc: number, p: any) => acc + p.balance, 0);

  const unpaidBills = properties.flatMap((p: any) => 
    p.bills.filter((b: any) => b.status === "unpaid")
      .map((b: any) => ({ ...b, propertyName: p.name, tenantName: p.tenant?.name || "Жилец" }))
  );

  return (
    <div className="flex-1 pb-24 px-4 py-6 md:px-8 max-w-4xl mx-auto w-full">
      <div className="mb-6">
        <p className="text-sm text-muted-foreground font-semibold tracking-wider">{currentTime}</p>
        <h1 className="text-3xl font-bold mt-1 text-slate-900 dark:text-white">С возвращением, {user.name}!</h1>
        <p className="text-sm mt-2 text-slate-700 dark:text-slate-300 italic">{compliment}</p>
      </div>

      {isTenant && (
        <>
          {myProperty && myProperty.tenant ? (
            <>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <Card className="col-span-1 md:col-span-2 border-slate-200">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-lg">Моя квартира</CardTitle>
                    <CardDescription>
                      {myProperty.city}, ул. {myProperty.street}, д. {myProperty.house}, кв. {myProperty.apartment}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-3 text-sm">
                    <div className="flex justify-between py-1 border-b border-dashed">
                      <span className="text-muted-foreground">Ежемесячная плата:</span>
                      <span className="font-semibold text-foreground">{myProperty.tenant.price.toLocaleString("ru-RU")} ₽</span>
                    </div>
                    <div className="flex justify-between py-1 border-b border-dashed">
                      <span className="text-muted-foreground">День оплаты:</span>
                      <span className="font-semibold">{myProperty.tenant.payment_day} числа каждого месяца</span>
                    </div>
                    <div className="flex justify-between py-1 border-b border-dashed">
                      <span className="text-muted-foreground">Передача показаний:</span>
                      <span className="font-semibold">{myProperty.tenant.reading_day} числа каждого месяца</span>
                    </div>
                    <div className="flex justify-between py-1 border-b border-dashed">
                      <span className="text-muted-foreground">Срок аренды:</span>
                      <span className="font-medium text-slate-600 dark:text-slate-300">
                        {new Date(myProperty.tenant.start_date).toLocaleDateString("ru-RU")} — {new Date(myProperty.tenant.end_date).toLocaleDateString("ru-RU")}
                      </span>
                    </div>
                  </CardContent>
                </Card>

                <Card className={`border ${myProperty.balance < 0 ? "border-destructive/30 bg-destructive/5" : "border-emerald-500/20 bg-emerald-500/5"}`}>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Ваш Баланс</CardTitle>
                  </CardHeader>
                  <CardContent className="flex flex-col justify-between h-32">
                    <div>
                      <h2 className={`text-3xl font-bold tracking-tight ${myProperty.balance < 0 ? "text-destructive" : "text-emerald-600"}`}>
                        {myProperty.balance.toLocaleString("ru-RU")} ₽
                      </h2>
                      <p className="text-xs text-muted-foreground mt-1">
                        {myProperty.balance < 0 
                          ? "У вас отрицательный баланс. Пожалуйста, оплатите счета!" 
                          : "Баланс в порядке. Спасибо!"}
                      </p>
                    </div>
                    <div className="flex gap-2 mt-2">
                      <Button 
                        onClick={() => setShowReadingModal(true)}
                        variant="outline"
                        className="flex-1 h-9 rounded-xl text-xs font-medium"
                      >
                        Показания
                      </Button>
                      <Button 
                        onClick={() => setShowTopUpModal(true)}
                        className="flex-1 h-9 rounded-xl text-xs font-medium"
                      >
                        Пополнить
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </div>

              <div className="mb-8">
                <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider mb-3">Расположение на карте</h3>
                <MapYandex properties={[myProperty]} className="w-full h-64 rounded-3xl overflow-hidden border border-slate-200 shadow-xs" />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
                <div>
                  <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                    <span>Мои счета</span>
                    {myProperty.bills.filter((b: any) => b.status === "unpaid").length > 0 && (
                      <Badge variant="destructive" className="rounded-full px-2 py-0.5 text-[10px]">
                        {myProperty.bills.filter((b: any) => b.status === "unpaid").length}
                      </Badge>
                    )}
                  </h3>
                  {myProperty.bills.length === 0 ? (
                    <p className="text-sm text-muted-foreground italic">Счетов пока нет.</p>
                  ) : (
                    <div className="space-y-4">
                      {[...myProperty.bills].map((bill) => (
                        <Card key={bill.id} className={`border ${bill.status === "paid" ? "border-border/50 bg-muted/20" : "border-destructive/30 bg-destructive/5"}`}>
                          <CardHeader className="p-4 pb-2">
                            <div className="flex justify-between items-center">
                              <span className="text-xs text-muted-foreground">
                                {new Date(bill.date).toLocaleDateString("ru-RU")}
                              </span>
                              <Badge variant={bill.status === "paid" ? "secondary" : "destructive"} className="text-xs font-semibold px-2.5 py-0.5">
                                {bill.status === "paid" ? "Оплачен" : "Ожидает оплаты"}
                              </Badge>
                            </div>
                            <CardTitle className="text-base mt-1 font-bold">
                              Счет на {bill.total.toLocaleString("ru-RU")} ₽
                            </CardTitle>
                          </CardHeader>
                          <CardContent className="p-4 pt-0 space-y-2 text-xs">
                            <div className="border-t pt-2 space-y-1">
                              {bill.items.map((item: any, idx: number) => (
                                <div key={idx} className="flex justify-between">
                                  <span className="text-muted-foreground">• {item.description}</span>
                                  <span className="font-medium">{item.amount.toLocaleString("ru-RU")} ₽</span>
                                </div>
                              ))}
                            </div>
                            <div className="flex justify-between text-muted-foreground border-t pt-2">
                              <span>Оплатить до:</span>
                              <span>{new Date(bill.due_date).toLocaleDateString("ru-RU")}</span>
                            </div>
                            <div className="flex gap-2 mt-3">
                              {bill.status === "unpaid" && (
                                <Button 
                                  onClick={() => handlePayBill(myProperty.id, bill.id, bill.total)}
                                  className="flex-1 h-9 rounded-xl text-xs font-medium"
                                >
                                  Оплатить ({bill.total} ₽)
                                </Button>
                              )}
                              <Button 
                                type="button"
                                onClick={() => {
                                  setSelectedInvoiceBill(bill);
                                  setShowInvoiceModal(true);
                                }}
                                variant="outline"
                                className={`h-9 rounded-xl text-xs font-medium ${bill.status === "paid" ? "w-full" : "px-3"}`}
                              >
                                Печать квитанции
                              </Button>
                            </div>
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  )}
                </div>

                <div>
                  <h3 className="text-lg font-bold mb-4">Переданные показания</h3>
                  {myProperty.readings.length === 0 ? (
                    <p className="text-sm text-muted-foreground italic">Показания пока не передавались.</p>
                  ) : (
                    <div className="space-y-4">
                      {[...myProperty.readings].map((reading: any) => (
                        <Card key={reading.id} className="border border-border/50 shadow-xs">
                          <CardHeader className="p-4 pb-2">
                            <div className="flex justify-between items-center">
                              <span className="text-xs text-muted-foreground">
                                {new Date(reading.date).toLocaleDateString("ru-RU")}
                              </span>
                              <Badge variant={reading.is_accounted ? "outline" : "outline"} className={reading.is_accounted ? "" : "border-amber-500/50 text-amber-600 bg-amber-500/5 animate-pulse"}>
                                {reading.is_accounted ? "Учтено в счете" : "Новое"}
                              </Badge>
                            </div>
                          </CardHeader>
                          <CardContent className="p-4 pt-0 grid grid-cols-3 gap-2 text-center text-xs">
                            <div className="bg-muted/40 p-2 rounded-xl">
                              <p className="text-[10px] text-muted-foreground">ГВС</p>
                              <p className="font-bold text-foreground">{reading.gvs} м³</p>
                            </div>
                            <div className="bg-muted/40 p-2 rounded-xl">
                              <p className="text-[10px] text-muted-foreground">ХВС</p>
                              <p className="font-bold text-foreground">{reading.hvs} м³</p>
                            </div>
                            <div className="bg-muted/40 p-2 rounded-xl">
                              <p className="text-[10px] text-muted-foreground">ЭЛ. Peak</p>
                              <p className="font-bold text-foreground">{reading.el1} кВт</p>
                            </div>
                            {reading.el2 !== null && reading.el2 !== undefined && (
                              <div className="bg-muted/40 p-2 rounded-xl col-span-3">
                                <p className="text-[10px] text-muted-foreground">ЭЛ. Off-Peak (T2)</p>
                                <p className="font-bold text-foreground">{reading.el2} кВт</p>
                              </div>
                            )}
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </>
          ) : (
            <Card className="mb-8 border-dashed border-border p-6 text-center">
              <h2 className="text-base font-bold text-foreground mt-3">У вас пока нет активной аренды</h2>
              <p className="text-xs text-muted-foreground max-w-sm mx-auto mt-1">
                Ознакомьтесь со списком предложений от арендодателей ниже, отправьте отклик, обсудите детали в чате и заселяйтесь!
              </p>
            </Card>
          )}

          <div className="mb-12">
            <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
              <span>Свободные квартиры для аренды</span>
              <Badge variant="secondary" className="text-[10px]">{vacantProperties.length}</Badge>
            </h3>

            {vacantProperties.length > 0 && (
              <div className="mb-6">
                <MapYandex 
                  properties={vacantProperties} 
                  className="w-full h-80 rounded-3xl overflow-hidden border border-border shadow-xs" 
                  selectedPropertyId={selectedPropertyId}
                  onMarkerClick={(prop) => {
                    setSelectedPropertyId(prop.id);
                    const el = document.getElementById(`vacant-prop-${prop.id}`);
                    if (el) {
                      el.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    }
                  }}
                />
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {vacantProperties.map((prop: any) => {
                const hasApplied = prop.applications?.some((a: any) => a.tenant_user_id === user.id);
                const isSelected = selectedPropertyId === prop.id;

                return (
                  <Card 
                    id={`vacant-prop-${prop.id}`}
                    key={prop.id} 
                    className={`border transition-all duration-300 ${isSelected ? "border-primary ring-2 ring-primary/20 shadow-md scale-[1.01]" : "border-border"}`}
                  >
                    <CardHeader className="pb-2">
                      <CardTitle className="text-base font-bold text-foreground leading-tight">
                        {prop.name}
                      </CardTitle>
                      <CardDescription className="text-xs">
                        {prop.city}, ул. {prop.street}, д. {prop.house}, кв. {prop.apartment}
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-3 text-xs">
                      <div className="flex justify-between py-1 border-b border-dashed">
                        <span className="text-muted-foreground">Рента квартиры:</span>
                        <span className="font-bold text-foreground text-sm">По договоренности</span>
                      </div>
                      <div className="grid grid-cols-3 gap-2 text-center text-[10px] bg-muted/40 p-2 rounded-xl border border-border/50">
                        <div>
                          <p className="text-muted-foreground">ГВС</p>
                          <p className="font-semibold text-foreground">{prop.gvs_tariff} ₽/м³</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">ХВС</p>
                          <p className="font-semibold text-foreground">{prop.hvs_tariff} ₽/м³</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">ЭЛ1</p>
                          <p className="font-semibold text-foreground">{prop.el1_tariff} ₽</p>
                        </div>
                      </div>

                      {hasApplied ? (
                        <Button 
                          disabled 
                          className="w-full text-muted-foreground font-semibold"
                        >
                          Вы откликнулись (В ожидании)
                        </Button>
                      ) : (
                        <Button 
                          onClick={() => handleApplyProperty(prop.id)}
                          className="w-full font-bold"
                        >
                          Откликнуться на предложение
                        </Button>
                      )}
                    </CardContent>
                  </Card>
                );
              })}

              {vacantProperties.length === 0 && (
                <p className="col-span-full text-center text-sm text-muted-foreground italic py-8 border rounded-2xl border-dashed">
                  Свободных объектов для аренды на данный момент нет.
                </p>
              )}
            </div>
          </div>
        </>
      )}

      {!isTenant && (
        <>
        <div className="p-2 border rounded-2xl px-4 mb-2">
          <p>Общий баланс:</p>
          <span className="font-semibold text-3xl">{totalSaldo} ₽</span>
        </div>
          {properties.length > 0 && (
            <div className="mb-8">
              <h3 className="text-lg font-bold mb-2">Мои объекты на карте</h3>
              <div className="flex justify-start items-center mb-2 gap-2">
                <div className="bg-muted px-4 py-1 rounded-full border">
                  <p><span className="text-xs">Всего:</span>  {totalProperties}</p>
                </div>
                <div className="bg-green-500/10 px-4 py-1 rounded-full border">
                  <p><span className="text-xs">Занято:</span> {occupiedCount}</p>
                </div>
                <div className="bg-yellow-500/10 px-4 py-1 rounded-full border">
                  <p><span className="text-xs">Свободно:</span> {vacantCount}</p>
                </div>
              </div>
              <MapYandex 
                properties={properties} 
                className="w-full h-50 rounded-3xl overflow-hidden border border-slate-200 shadow-xs"
                selectedPropertyId={selectedPropertyId}
                onMarkerClick={(prop) => setSelectedPropertyId(prop.id)}
              />
            </div>
          )}
          <div className="mb-8">
            <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
              <span>Горящие события</span>
              {totalUrgentCount > 0 && (
                <Badge variant="destructive" className="rounded-full px-1.5 py-0.5 text-xs bg-red-500 text-white">
                  {totalUrgentCount}
                </Badge>
              )}
            </h3>

            {totalUrgentCount === 0 ? (
              <Card className="border-dashed border-border bg-muted/20 p-6 text-center">
                <p className="text-sm text-muted-foreground font-semibold flex items-center justify-center gap-2">
                  <span>Все идет по плану! Нет событий, требующих вашего внимания.</span>
                </p>
              </Card>
            ) : (
              <div className="space-y-3">
                {negBalanceProps.map(p => (
                  <div key={`neg-${p.id}`} className="flex items-start gap-3 bg-destructive/5 border border-destructive/20 rounded-2xl p-4">
                    <div className="flex-1 text-xs">
                      <p className="font-bold text-destructive">Отрицательный баланс объекта:</p>
                      <p className="text-muted-foreground mt-0.5">
                        У объекта <strong className="font-semibold text-foreground">"{p.name}"</strong> баланс составляет <strong className="text-destructive font-bold">{p.balance.toLocaleString("ru-RU")} ₽</strong>.
                      </p>
                    </div>
                    <Button size="sm" variant="outline" className="h-7 text-[10px] shrink-0" onClick={() => router.push(`/services`)}>
                      Подробнее
                    </Button>
                  </div>
                ))}

                {checkoutProps.map(p => {
                  const end = new Date(p.tenant.end_date);
                  const daysLeft = Math.ceil((end.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));
                  const isExpired = daysLeft < 0;

                  return (
                    <div key={`checkout-${p.id}`} className="flex items-start gap-3 bg-amber-500/5 border border-amber-500/20 rounded-2xl p-4">
                      <div className="flex-1 text-xs">
                        <p className="font-bold text-amber-600">Предупреждение о выезде арендатора:</p>
                        <p className="text-muted-foreground mt-0.5">
                          {isExpired ? (
                            <>
                              Срок аренды арендатора <strong className="font-semibold text-foreground">{p.tenant.name}</strong> в <strong className="font-semibold">"{p.name}"</strong> истек <strong className="font-bold text-amber-600">{Math.abs(daysLeft)} дн. назад</strong>!
                            </>
                          ) : (
                            <>
                              Арендатор <strong className="font-semibold text-foreground">{p.tenant.name}</strong> съезжает из <strong className="font-semibold">"{p.name}"</strong> через <strong className="font-bold text-amber-600">{daysLeft} дн.</strong> ({end.toLocaleDateString("ru-RU")}).
                            </>
                          )}
                        </p>
                      </div>
                      <Button size="sm" variant="outline" className="h-7 text-[10px] shrink-0" onClick={() => router.push(`/services`)}>
                        Управлять
                      </Button>
                    </div>
                  );
                })}

                {readingsProps.map(p => (
                  <div key={`readings-${p.id}`} className="flex items-start gap-3 bg-primary/5 border border-primary/20 rounded-2xl p-4">
                    <div className="flex-1 text-xs">
                      <p className="font-bold text-primary">Показания не переданы:</p>
                      <p className="text-muted-foreground mt-0.5">
                        До дня оплаты ({p.tenant.payment_day} числа) по <strong className="font-semibold text-foreground">"{p.name}"</strong> осталось менее 5 дней. Показания за текущий цикл не передавались!
                      </p>
                    </div>
                    <Button size="sm" variant="outline" className="h-7 text-[10px] shrink-0" onClick={() => router.push(`/services`)}>
                      Напомнить
                    </Button>
                  </div>
                ))}

                {recentBills.map(b => (
                  <div key={`bill-gen-${b.id}`} className="flex items-start gap-3 bg-emerald-500/5 border border-emerald-500/20 rounded-2xl p-4">
                    <div className="flex-1 text-xs">
                      <p className="font-semibold text-emerald-900 text-sm">Модуль генерации счетов:</p>
                      <p className="text-muted-foreground mt-0.5">
                        Авто-сгенерирован счет для <strong className="font-semibold text-foreground">"{b.propertyName}"</strong> от {new Date(b.date).toLocaleDateString("ru-RU")} на сумму <strong className="font-bold text-emerald-600">{b.total.toLocaleString("ru-RU")} ₽</strong>.
                      </p>
                    </div>
                    <Button size="sm" variant="outline" className="h-7 text-[12px] shrink-0" onClick={() => router.push(`/services`)}>
                      История
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold">Недвижимость</h3>
              <Button
                variant="outline"
                size="sm"
                onClick={() => router.push("/services")}
                className="text-xs rounded-xl"
              >
                Подробнее
              </Button>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {properties.map((prop: any) => {
                const isOccupied = !!prop.tenant;
                return (
                  <Card key={prop.id} className="hover:shadow-md transition-all border-border duration-200 bg-card">
                    <CardHeader className="pb-2">
                      <div className="flex justify-between items-start">
                        <Badge variant={isOccupied ? "secondary" : "destructive"}>
                          {isOccupied ? "Занято" : "Свободно"}
                        </Badge>
                      </div>
                      <CardTitle className="text-base mt-1 font-bold text-foreground leading-tight">
                        {prop.name}
                      </CardTitle>
                      <CardDescription className="text-xs">
                        {prop.city}, ул. {prop.street}, д. {prop.house}
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-2 text-xs">
                      <div className="flex justify-between items-center py-1 border-b border-dashed border-border">
                        <span className="text-muted-foreground">Баланс:</span>
                        <span className={`font-bold text-sm ${prop.balance < 0 ? "text-destructive" : "text-emerald-600"}`}>
                          {prop.balance.toLocaleString("ru-RU")} ₽
                        </span>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}

              {properties.length === 0 && (
                <div className="col-span-full py-8 text-center">
                  <p className="text-sm text-muted-foreground italic">У вас еще нет добавленных объектов.</p>
                </div>
              )}
            </div>
          </div>

          {unpaidBills.length > 0 && (
            <div className="mb-8">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2 text-destructive">
                <span>Ожидают оплаты жильцами</span>
                <Badge variant="destructive" className="rounded-full px-2 py-0.5 font-bold text-xs">
                  {unpaidBills.length}
                </Badge>
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {unpaidBills.map((bill: any) => (
                  <Card key={bill.id} className="border border-destructive/20 bg-destructive/5 shadow-xs">
                    <CardHeader className="p-4 pb-2 flex flex-row items-center justify-between">
                      <div>
                        <CardTitle className="text-sm font-bold text-foreground">{bill.propertyName}</CardTitle>
                        <CardDescription className="text-[10px]">Арендатор: {bill.tenantName}</CardDescription>
                      </div>
                      <span className="text-destructive font-bold text-sm">{bill.total.toLocaleString("ru-RU")} ₽</span>
                    </CardHeader>
                    <CardContent className="p-4 pt-0 text-[10px] space-y-2">
                      <div className="border-t border-dashed border-border/50 pt-2">
                        {bill.items.map((item: any, idx: number) => (
                          <div key={idx} className="flex justify-between text-muted-foreground">
                            <span>• {item.description}</span>
                            <span className="font-semibold text-foreground">{item.amount.toLocaleString("ru-RU")} ₽</span>
                          </div>
                        ))}
                      </div>
                      <div className="flex justify-between text-muted-foreground border-t border-dashed border-border/50 pt-2">
                        <span>Выставлен: {new Date(bill.date).toLocaleDateString("ru-RU")}</span>
                        <span className="font-bold text-destructive">Срок: {new Date(bill.due_date).toLocaleDateString("ru-RU")}</span>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          )}

        </>
      )}

      <Dialog open={showReadingModal} onOpenChange={setShowReadingModal}>
        <DialogContent className="w-full max-w-sm">
          <DialogHeader>
            <DialogTitle>Передать показания ЖКХ</DialogTitle>
            <DialogDescription>Введите текущие значения счетчиков</DialogDescription>
          </DialogHeader>
          <form onSubmit={handleAddReading} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="gvs">Горячее водоснабжение (ГВС) *</Label>
              <Input 
                id="gvs" 
                type="number" 
                step="0.01" 
                placeholder="Например: 124.5" 
                value={gvs} 
                onChange={e => setGvs(e.target.value)} 
                required 
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="hvs">Холодное водоснабжение (ХВС) *</Label>
              <Input 
                id="hvs" 
                type="number" 
                step="0.01" 
                placeholder="Например: 255.2" 
                value={hvs} 
                onChange={e => setHvs(e.target.value)} 
                required 
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="el1">Электричество Т1 (Пик) *</Label>
              <Input 
                id="el1" 
                type="number" 
                step="0.01" 
                placeholder="Например: 1540" 
                value={el1} 
                onChange={e => setEl1(e.target.value)} 
                required 
              />
            </div>
            {myProperty && myProperty.el2_tariff !== null && (
              <div className="space-y-2">
                <Label htmlFor="el2">Электричество Т2 (Ночь)</Label>
                <Input 
                  id="el2" 
                  type="number" 
                  step="0.01" 
                  placeholder="Например: 812" 
                  value={el2} 
                  onChange={e => setEl2(e.target.value)} 
                />
              </div>
            )}
            <div className="flex gap-2 justify-end pt-2">
              <Button type="button" variant="outline" onClick={() => setShowReadingModal(false)}>
                Отмена
              </Button>
              <Button type="submit">
                Передать
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      
      {myProperty && (
        <Dialog open={showTopUpModal} onOpenChange={setShowTopUpModal}>
          <DialogContent className="w-full max-w-sm">
            <DialogHeader>
              <DialogTitle>Пополнить баланс</DialogTitle>
              <DialogDescription>Внесите произвольную сумму для оплаты аренды или коммунальных платежей</DialogDescription>
            </DialogHeader>
            <form onSubmit={handleTopUpSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="topUpAmt">Сумма пополнения (₽) *</Label>
                <Input 
                  id="topUpAmt" 
                  type="number" 
                  placeholder="Например: 5000" 
                  value={topUpAmount} 
                  onChange={e => setTopUpAmount(e.target.value)} 
                  required 
                />
              </div>
              <div className="flex gap-2 justify-end pt-2">
                <Button type="button" variant="outline" onClick={() => setShowTopUpModal(false)}>
                  Отмена
                </Button>
                <Button type="submit">
                  Пополнить
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      )}

      
      {myProperty && selectedInvoiceBill && (
        <Dialog open={showInvoiceModal} onOpenChange={(open) => {
          setShowInvoiceModal(open);
          if (!open) setSelectedInvoiceBill(null);
        }}>
          <DialogContent className="w-full max-w-xl max-h-[90dvh] flex flex-col" showCloseButton>
            <DialogHeader className="border-b pb-3 shrink-0">
              <DialogTitle>Квитанция на оплату ЖКХ</DialogTitle>
              <DialogDescription>Счет № {selectedInvoiceBill.id} от {new Date(selectedInvoiceBill.date).toLocaleDateString("ru-RU")}</DialogDescription>
            </DialogHeader>
            <div className="overflow-y-auto p-6 space-y-6 flex-1">
              <div className="border border-border p-8 rounded-xl bg-card text-card-foreground text-xs shadow-sm space-y-4 font-mono leading-relaxed">
                <h2 className="text-center text-sm font-bold border-b-2 border-border pb-2 uppercase">СЧЕТ НА ОПЛАТУ АРЕНДЫ И УСЛУГ ЖКХ</h2>
                
                <div className="grid grid-cols-2 gap-4 text-[11px]">
                  <div>
                    <p className="font-bold">Получатель (Арендодатель):</p>
                    <p>{myProperty.landlordName || "Владелец"}</p>
                    <p>Телефон: {myProperty.landlordPhone || "Уточните телефон"}</p>
                  </div>
                  <div>
                    <p className="font-bold">Плательщик (Арендатор):</p>
                    <p>{user.name}</p>
                    <p>Объект: {myProperty.city}, ул. {myProperty.street}, д. {myProperty.house}, кв. {myProperty.apartment}</p>
                  </div>
                </div>

                <div className="border-t border-b py-2 my-2 border-border">
                  <table className="w-full text-[10px]">
                    <thead>
                      <tr className="border-b text-left border-border">
                        <th className="py-1">Описание услуги</th>
                        <th className="py-1 text-right">Сумма</th>
                      </tr>
                    </thead>
                    <tbody>
                      {selectedInvoiceBill.items.map((item: any, idx: number) => (
                        <tr key={idx} className="border-b border-dashed border-border">
                          <td className="py-1">{item.description}</td>
                          <td className="py-1 text-right">{item.amount.toLocaleString("ru-RU")} ₽</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <div className="text-right text-xs">
                  <p className="font-bold text-sm">Итого к оплате: {selectedInvoiceBill.total.toLocaleString("ru-RU")} ₽</p>
                  <p className="text-[10px] text-muted-foreground mt-1">Оплатить в срок до: {new Date(selectedInvoiceBill.due_date).toLocaleDateString("ru-RU")}</p>
                </div>
              </div>

              <div className="flex justify-end gap-2 pt-2 border-t border-border">
                <Button variant="outline" onClick={() => {
                  setShowInvoiceModal(false);
                  setSelectedInvoiceBill(null);
                }}>
                  Закрыть
                </Button>
                <Button onClick={() => window.print()}>
                  Печать
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
