"use client";

import { useEffect, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/use-user";
import { useQueryClient } from "@tanstack/react-query";
import { useProperties, useVacantProperties, useUnlinkedTenants } from "@/hooks/use-properties";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { AuthGuard } from "@/components/auth-guard";
import PropertyDetailSidebar from "@/components/property-detail-sidebar";
import PropertyFormDialog from "@/components/property-form-dialog";
import LeaseFormDialog from "@/components/lease-form-dialog";
import { HugeiconsIcon } from "@hugeicons/react";
import { Edit03Icon, Delete02Icon } from "@hugeicons/core-free-icons";
import { 
  getProperties, 
  createProperty, 
  updateProperty, 
  deleteProperty, 
  createLease, 
  deleteLease, 
  submitReadings, 
  payBill, 
  addCustomItem,
  getVacantProperties,
  getUnlinkedTenants
} from "@/lib/properties";

export default function ServicesPage() {
  return (
    <AuthGuard mode="auth">
      <ServicesContent />
    </AuthGuard>
  );
}

function ServicesContent() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: user, isLoading } = useUser();

  const [selectedPropertyId, setSelectedPropertyId] = useState<string | null>(null);

  const [showAddModal, setShowAddModal] = useState(false);
  const [editingProperty, setEditingProperty] = useState<any | null>(null);
  const [showTenantModal, setShowTenantModal] = useState(false);
  const [showContractModal, setShowContractModal] = useState(false);
  const [showInvoiceModal, setShowInvoiceModal] = useState(false);
  const [selectedBill, setSelectedBill] = useState<any | null>(null);

  
  const [propName, setPropName] = useState("");
  const [coords, setCoords] = useState("");
  const [country, setCountry] = useState("");
  const [region, setRegion] = useState("");
  const [city, setCity] = useState("");
  const [street, setStreet] = useState("");
  const [house, setHouse] = useState("");
  const [apartment, setApartment] = useState("");
  const [gvsTariff, setGvsTariff] = useState("");
  const [hvsTariff, setHvsTariff] = useState("");
  const [el1Tariff, setEl1Tariff] = useState("");
  const [el2Tariff, setEl2Tariff] = useState("");
  const [formError, setFormError] = useState("");

  const [selectedTenantUserId, setSelectedTenantUserId] = useState("");
  const [tenantName, setTenantName] = useState("");
  const [tenantDoc, setTenantDoc] = useState("");
  const [tenantPhone, setTenantPhone] = useState("");
  const [monthsOfRent, setMonthsOfRent] = useState("");
  const [rentPrice, setRentPrice] = useState("");
  const [paymentDay, setPaymentDay] = useState("5");
  const [readingDay, setReadingDay] = useState("1");
  const [tenantError, setTenantError] = useState("");

  const [gvsVal, setGvsVal] = useState("");
  const [hvsVal, setHvsVal] = useState("");
  const [el1Val, setEl1Val] = useState("");
  const [el2Val, setEl2Val] = useState("");

  const [paymentAmount, setPaymentAmount] = useState("");
  const [customDesc, setCustomDesc] = useState("");
  const [customAmount, setCustomAmount] = useState("");

  const isLandlord = !!user && user.role === "landlord";
  const { data: rawPropertiesData } = useProperties({ enabled: isLandlord });
  const { data: vacantData } = useVacantProperties({ enabled: isLandlord });
  const { data: unlinkedTenantsData } = useUnlinkedTenants({ enabled: isLandlord });

  const unlinkedTenants = unlinkedTenantsData || [];

  const properties = useMemo(() => {
    if (!rawPropertiesData) return [];
    const vacantList = vacantData || [];
    const tenantsList = unlinkedTenantsData || [];

    return rawPropertiesData.map((prop: any) => {
      const matchingVacant = vacantList.find((v: any) => v.id === prop.id);
      const rawApps = matchingVacant?.applications || [];
      
      const enrichedApps = rawApps.map((app: any) => {
        const tenantInfo = tenantsList.find((t: any) => t.id === app.tenant_user_id);
        return {
          ...app,
          name: tenantInfo?.name || app.name || "Неизвестный жилец",
          phone: tenantInfo?.phone || app.phone || "Не указан",
          document: tenantInfo?.document || app.document || "Не указан",
        };
      });

      return {
        ...prop,
        applications: enrichedApps,
      };
    });
  }, [rawPropertiesData, vacantData, unlinkedTenantsData]);

  const selectedProperty = useMemo(() => {
    if (!selectedPropertyId) return null;
    return properties.find((p: any) => p.id === selectedPropertyId) || null;
  }, [properties, selectedPropertyId]);

  const setSelectedProperty = (prop: any | null) => {
    setSelectedPropertyId(prop ? prop.id : null);
  };

  const loadProperties = async () => {
    queryClient.invalidateQueries({ queryKey: ["properties"] });
    queryClient.invalidateQueries({ queryKey: ["vacantProperties"] });
    queryClient.invalidateQueries({ queryKey: ["unlinkedTenants"] });
  };

  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background text-foreground">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
      </div>
    );
  }

  if (!user || user.role !== "landlord") {
    return (
      <div className="p-6 text-center max-w-lg mx-auto mt-20">
        <Card className="border-rose-300">
          <CardHeader>
            <CardTitle className="text-rose-600">Доступ ограничен</CardTitle>
            <CardDescription>Управление объектами доступно только Арендодателям.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => router.push("/home")}>На главную</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const handleOpenAddModal = (p: any = null) => {
    if (p) {
      setEditingProperty(p);
      setPropName(p.name);
      setCoords(p.coordinates);
      setCountry(p.country || "");
      setRegion(p.region);
      setCity(p.city);
      setStreet(p.street);
      setHouse(p.house);
      setApartment(p.apartment);
      setGvsTariff(p.gvs_tariff.toString());
      setHvsTariff(p.hvs_tariff.toString());
      setEl1Tariff(p.el1_tariff.toString());
      setEl2Tariff(p.el2_tariff?.toString() || "");
    } else {
      setEditingProperty(null);
      setPropName("");
      setCoords("55.7558, 37.6173");
      setCountry("Россия");
      setRegion("Московская область");
      setCity("Москва");
      setStreet("");
      setHouse("");
      setApartment("");
      setGvsTariff("220");
      setHvsTariff("50");
      setEl1Tariff("6.5");
      setEl2Tariff("");
    }
    setFormError("");
    setShowAddModal(true);
  };

  const handlePropertySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError("");

    const coordsRegex = /^-?\d+(\.\d+)?,\s*-?\d+(\.\d+)?$/;
    if (!coordsRegex.test(coords)) {
      setFormError("Координаты должны быть в формате 'широта, долгота' (например: 55.7558, 37.6173)");
      return;
    }

    const gvsNum = parseFloat(gvsTariff);
    const hvsNum = parseFloat(hvsTariff);
    const el1Num = parseFloat(el1Tariff);
    if (isNaN(gvsNum) || gvsNum < 0 || isNaN(hvsNum) || hvsNum < 0 || isNaN(el1Num) || el1Num < 0) {
      setFormError("Тарифы ЖКХ должны быть положительными числами.");
      return;
    }

    const compiledName = propName.trim() || `${street} ${Math.floor(Math.random() * 99) + 1}`;

    try {
      const payload = {
        name: compiledName,
        coordinates: coords,
        country,
        region,
        city,
        street,
        house,
        apartment,
        gvsTariff: gvsNum,
        hvsTariff: hvsNum,
        el1Tariff: el1Num,
        el2Tariff: el2Tariff ? parseFloat(el2Tariff) : null,
      };

      if (editingProperty) {
        await updateProperty(editingProperty.id, payload);
      } else {
        await createProperty(payload);
      }

      setShowAddModal(false);
      loadProperties();
    } catch (err: any) {
      setFormError(err.message || "Ошибка сохранения");
    }
  };

  const handlePropertyDelete = async (id: string, name: string) => {
    if (!confirm(`Вы действительно хотите удалить объект "${name}"?`)) return;

    try {
      await deleteProperty(id);
      setSelectedProperty(null);
      loadProperties();
    } catch (err: any) {
      alert(err.message || "Ошибка удаления объекта");
    }
  };

  const handleAddTenantSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTenantError("");

    if (!selectedTenantUserId) {
      setTenantError("Пожалуйста, выберите зарегистрированного жильца.");
      return;
    }

    const priceNum = parseFloat(rentPrice);
    const monthsNum = parseInt(monthsOfRent);

    if (isNaN(priceNum) || priceNum <= 0 || isNaN(monthsNum) || monthsNum <= 0) {
      setTenantError("Срок аренды и цена должны быть положительными числами.");
      return;
    }

    try {
      await createLease(selectedProperty.id, {
        tenantUserId: selectedTenantUserId,
        price: priceNum,
        monthsOfRent: monthsNum,
        paymentDay: parseInt(paymentDay),
        readingDay: parseInt(readingDay),
      });

      setShowTenantModal(false);
      setSelectedTenantUserId("");
      setTenantName("");
      setTenantDoc("");
      setTenantPhone("");
      loadProperties();
    } catch (err: any) {
      setTenantError(err.message || "Ошибка оформления");
    }
  };

  const handleRemoveTenant = async (propId: string) => {
    if (!confirm("Выселить арендатора?")) return;

    try {
      await deleteLease(propId);
      loadProperties();
    } catch (err: any) {
      alert(err.message || "Ошибка выселения");
    }
  };

  const handleAddReadingSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProperty) return;

    try {
      await submitReadings(selectedProperty.id, {
        gvs: parseFloat(gvsVal),
        hvs: parseFloat(hvsVal),
        el1: parseFloat(el1Val),
        el2: el2Val ? parseFloat(el2Val) : null,
      });

      setGvsVal("");
      setHvsVal("");
      setEl1Val("");
      setEl2Val("");
      loadPropertiesAndSelected();
    } catch (err: any) {
      alert(err.message || "Ошибка отправки показаний");
    }
  };

  const handleLogPaymentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProperty) return;

    try {
      await payBill(selectedProperty.id, { amount: parseFloat(paymentAmount) });
      setPaymentAmount("");
      loadPropertiesAndSelected();
    } catch (err: any) {
      alert(err.message || "Ошибка регистрации оплаты");
    }
  };

  const handleAddCustomItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProperty || !customDesc || !customAmount) return;

    try {
      await addCustomItem(selectedProperty.id, {
        description: customDesc,
        amount: parseFloat(customAmount)
      });
      setCustomDesc("");
      setCustomAmount("");
      loadPropertiesAndSelected();
      alert("Элемент успешно добавлен и будет учтен в следующем авто-счете.");
    } catch (err: any) {
      alert(err.message || "Ошибка добавления услуги");
    }
  };

  const loadPropertiesAndSelected = async () => {
    await loadProperties();
  };

  const checkCurrentMonthReadings = (prop: any) => {
    if (prop.readings.length === 0) return false;
    const latest = prop.readings[0];
    const readingDate = new Date(latest.date);
    const curDate = new Date();
    return readingDate.getMonth() === curDate.getMonth() && readingDate.getFullYear() === curDate.getFullYear();
  };

  return (
    <div className="flex-1 pb-24 px-4 py-6 md:px-8 max-w-5xl mx-auto w-full">
      
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Недвижимость</h1>
          <p className="text-xs text-muted-foreground mt-0.5">Управление объектами аренды и арендаторами</p>
        </div>
        {properties.length > 0 && (
          <div className="flex gap-2">
            <Button 
              onClick={() => handleOpenAddModal()} 
            >
              Добавить объект
            </Button>
          </div>
        )}
      </div>

      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {properties.map((prop: any) => {
          const isOccupied = !!prop.tenant;
          const isReadingUpdated = checkCurrentMonthReadings(prop);

          return (
            <Card key={prop.id} className="hover:shadow-md transition-all border-border duration-200 bg-card">
              <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                  <Badge variant={isOccupied ? "secondary" : "destructive"}>
                    {isOccupied ? "Занято" : "Свободно"}
                  </Badge>
                  <div className="flex gap-1.5">
                    <Button 
                      variant="ghost" 
                      size="icon-lg"
                      className="size-9 bg-muted hover:bg-neutral-400/40"
                      onClick={() => handleOpenAddModal(prop)} 
                      title="Редактировать"
                    >
                      <HugeiconsIcon icon={Edit03Icon} className="size-5" />
                    </Button>
                    <Button 
                      variant="destructive" 
                      size="icon-lg"
                      className="size-9"
                      onClick={() => handlePropertyDelete(prop.id, prop.name)} 
                      title="Удалить"
                    >
                      <HugeiconsIcon icon={Delete02Icon} className="size-5" />
                    </Button>
                  </div>
                </div>
                <CardTitle className="text-lg mt-1 font-bold text-foreground leading-tight">
                  {prop.name}
                </CardTitle>
                <CardDescription className="text-xs">
                  {prop.city}, ул. {prop.street}, д. {prop.house}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3.5 text-xs">
                <div className="flex justify-between items-center py-1 border-b border-dashed border-border">
                  <span className="text-muted-foreground">Баланс объекта:</span>
                  <span className={`font-bold text-sm ${prop.balance < 0 ? "text-destructive" : "text-emerald-600"}`}>
                    {prop.balance.toLocaleString("ru-RU")} ₽
                  </span>
                </div>
                <div className="flex justify-between items-center py-1 border-b border-dashed border-border">
                  <span className="text-muted-foreground">Показания за месяц:</span>
                  <span className={`font-semibold ${isReadingUpdated ? "text-emerald-600" : "text-amber-600"}`}>
                    {isReadingUpdated ? "Переданы" : "Ожидаются"}
                  </span>
                </div>

                <Button 
                  onClick={() => setSelectedProperty(prop)} 
                  variant="outline"
                  className="w-full text-xs py-2 rounded-xl"
                >
                  Подробнее о квартире
                </Button>
              </CardContent>
            </Card>
          );
        })}

        {properties.length === 0 && (
          <div className="col-span-full py-16 text-center">
            <p className="text-lg text-muted-foreground italic">У вас еще нет добавленных объектов.</p>
            <Button onClick={() => handleOpenAddModal()} className="mt-4">Добавить первый объект</Button>
          </div>
        )}
      </div>

      
      {selectedProperty && (
        <PropertyDetailSidebar
          property={selectedProperty}
          user={user}
          onClose={() => setSelectedProperty(null)}
          onContractClick={() => setShowContractModal(true)}
          onRemoveTenant={handleRemoveTenant}
          onSubmitReadings={handleAddReadingSubmit}
          onLogPayment={handleLogPaymentSubmit}
          onAddCustomItem={handleAddCustomItem}
          onShowTenantModal={(app) => {
            setSelectedTenantUserId(app.tenant_user_id);
            setTenantName(app.name);
            setTenantDoc(app.document);
            setTenantPhone(app.phone);
            setShowTenantModal(true);
          }}
          onManualLeaseClick={() => {
            setSelectedTenantUserId("");
            setTenantName("");
            setTenantDoc("");
            setTenantPhone("");
            setShowTenantModal(true);
          }}
          onBillClick={(bill) => {
            setSelectedBill(bill);
            setShowInvoiceModal(true);
          }}
          gvsVal={gvsVal} setGvsVal={setGvsVal}
          hvsVal={hvsVal} setHvsVal={setHvsVal}
          el1Val={el1Val} setEl1Val={setEl1Val}
          el2Val={el2Val} setEl2Val={setEl2Val}
          paymentAmount={paymentAmount} setPaymentAmount={setPaymentAmount}
          customDesc={customDesc} setCustomDesc={setCustomDesc}
          customAmount={customAmount} setCustomAmount={setCustomAmount}
        />
      )}


      <PropertyFormDialog
        open={showAddModal}
        onOpenChange={setShowAddModal}
        editingProperty={editingProperty}
        propName={propName} setPropName={setPropName}
        coords={coords} setCoords={setCoords}
        country={country} setCountry={setCountry}
        region={region} setRegion={setRegion}
        city={city} setCity={setCity}
        street={street} setStreet={setStreet}
        house={house} setHouse={setHouse}
        apartment={apartment} setApartment={setApartment}
        gvsTariff={gvsTariff} setGvsTariff={setGvsTariff}
        hvsTariff={hvsTariff} setHvsTariff={setHvsTariff}
        el1Tariff={el1Tariff} setEl1Tariff={setEl1Tariff}
        el2Tariff={el2Tariff} setEl2Tariff={setEl2Tariff}
        formError={formError}
        onSubmit={handlePropertySubmit}
      />


      <LeaseFormDialog
        open={showTenantModal}
        onOpenChange={setShowTenantModal}
        propertyName={selectedProperty?.name || ""}
        tenantName={tenantName}
        tenantDoc={tenantDoc}
        tenantPhone={tenantPhone}
        monthsOfRent={monthsOfRent} setMonthsOfRent={setMonthsOfRent}
        rentPrice={rentPrice} setRentPrice={setRentPrice}
        paymentDay={paymentDay} setPaymentDay={setPaymentDay}
        readingDay={readingDay} setReadingDay={setReadingDay}
        tenantError={tenantError}
        unlinkedTenants={unlinkedTenants}
        selectedTenantUserId={selectedTenantUserId}
        setSelectedTenantUserId={setSelectedTenantUserId}
        onSubmit={handleAddTenantSubmit}
      />

      
      <Dialog open={showContractModal && !!selectedProperty?.tenant} onOpenChange={setShowContractModal}>
        <DialogContent className="w-full max-w-2xl max-h-[90dvh] flex flex-col" showCloseButton>
          <DialogHeader className="shrink-0 border-b pb-3">
            <DialogTitle>Печать договора аренды жилья</DialogTitle>
            <DialogDescription>Сгенерирован автоматически системой My Rent 2</DialogDescription>
          </DialogHeader>
          <div className="overflow-y-auto p-6 space-y-6">
            {selectedProperty?.tenant && (
              <div className="border border-slate-300 p-8 rounded-xl bg-white text-slate-800 text-xs shadow-sm max-h-96 overflow-y-auto leading-relaxed space-y-4 font-serif">
                <h2 className="text-center text-sm font-bold tracking-wider uppercase">ДОГОВОР АРЕНДЫ ЖИЛОГО ПОМЕЩЕНИЯ</h2>
                <div className="flex justify-between font-bold">
                  <span>г. {selectedProperty.city}</span>
                  <span>{new Date(selectedProperty.tenant.start_date).toLocaleDateString("ru-RU")} г.</span>
                </div>
                <p>
                  Гражданин(ка) <strong>{user.name}</strong>, документ: {user.document || "Паспорт РФ"}, именуемый(ая) в дальнейшем «Арендодатель», с одной стороны, и гражданин(ка) <strong>{selectedProperty.tenant.name}</strong>, паспорт: {selectedProperty.tenant.document}, именуемый(ая) в дальнейшем «Арендатор», с другой стороны, заключили настоящий договор о следующем:
                </p>
                <div>
                  <h4 className="font-bold border-b mb-1 uppercase text-[10px]">1. Предмет Договора</h4>
                  <p>
                    1.1. Арендодатель предоставляет Арендатору во временное возмездное владение и пользование жилое помещение, квартиру, расположенную по адресу: <strong>г. {selectedProperty.city}, ул. {selectedProperty.street}, д. {selectedProperty.house}, кв. {selectedProperty.apartment}</strong>.
                  </p>
                </div>
                <div>
                  <h4 className="font-bold border-b mb-1 uppercase text-[10px]">2. Срок Аренды</h4>
                  <p>
                    2.1. Срок аренды устанавливается на {selectedProperty.tenant.months_of_rent} месяцев (из расчета 30 дней в месяце) с <strong>{new Date(selectedProperty.tenant.start_date).toLocaleDateString("ru-RU")} г.</strong> по <strong>{new Date(selectedProperty.tenant.end_date).toLocaleDateString("ru-RU")} г.</strong>
                  </p>
                </div>
                <div>
                  <h4 className="font-bold border-b mb-1 uppercase text-[10px]">3. Платежи и расчеты</h4>
                  <p>
                    3.1. Плата за наем Жилого помещения составляет <strong>{selectedProperty.tenant.price.toLocaleString("ru-RU")} рублей</strong> в месяц.
                  </p>
                  <p>
                    3.2. Оплата производится ежемесячно не позднее <strong>{selectedProperty.tenant.payment_day} числа</strong> текущего месяца на банковскую карту Арендодателя.
                  </p>
                </div>
                <div>
                  <h4 className="font-bold border-b mb-1 uppercase text-[10px]">4. Подписи сторон</h4>
                  <div className="grid grid-cols-2 gap-4 mt-6 border-t pt-4">
                    <div>
                      <p className="font-bold">Арендодатель:</p>
                      <p className="mt-8">____________ / {user.name} /</p>
                    </div>
                    <div>
                      <p className="font-bold">Арендатор:</p>
                      <p className="mt-8">____________ / {selectedProperty.tenant.name} /</p>
                    </div>
                  </div>
                </div>
              </div>
            )}
            <div className="flex justify-end gap-2 pt-2 border-t">
              <Button variant="outline" onClick={() => setShowContractModal(false)}>Закрыть</Button>
              <Button onClick={() => window.print()} className="bg-slate-800 text-white font-bold">Напечатать</Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>


      <Dialog open={showInvoiceModal && !!selectedBill} onOpenChange={(open) => { setShowInvoiceModal(open); if (!open) setSelectedBill(null); }}>
        <DialogContent className="w-full max-w-xl max-h-[90dvh] flex flex-col" showCloseButton>
          <DialogHeader className="shrink-0 border-b pb-3">
            <DialogTitle>Печать счета на оплату ЖКХ</DialogTitle>
            <DialogDescription>Счет № {selectedBill?.id} от {selectedBill?.date ? new Date(selectedBill?.date).toLocaleDateString("ru-RU") : ""}</DialogDescription>
          </DialogHeader>
          <div className="overflow-y-auto p-6 space-y-6">
            {selectedBill && selectedProperty && (
              <>
                <div className="border border-slate-300 p-8 rounded-xl bg-white text-slate-800 text-xs shadow-sm space-y-4 font-mono leading-relaxed">
                  <h2 className="text-center text-sm font-bold border-b-2 border-slate-950 pb-2 uppercase">СЧЕТ НА ОПЛАТУ АРЕНДЫ И УСЛУГ ЖКХ</h2>
                  <div className="grid grid-cols-2 gap-4 text-[11px]">
                    <div>
                      <p className="font-bold">Исполнитель:</p>
                      <p>{user.name}</p>
                      <p>Реквизиты: {user.paymentCard || "Уточните карту для перевода"}</p>
                      <p>Телефон: {user.phone}</p>
                    </div>
                    <div>
                      <p className="font-bold">Плательщик (Арендатор):</p>
                      <p>{selectedProperty.tenant?.name || "Жилец"}</p>
                      <p>Объект: {selectedProperty.city}, ул. {selectedProperty.street}, д. {selectedProperty.house}, кв. {selectedProperty.apartment}</p>
                    </div>
                  </div>
                  <div className="border-t border-b py-2 my-2">
                    <table className="w-full text-[10px]">
                      <thead>
                        <tr className="border-b text-left">
                          <th className="py-1">Описание услуги</th>
                          <th className="py-1 text-right">Сумма к оплате</th>
                        </tr>
                      </thead>
                      <tbody>
                        {selectedBill.items.map((item: any, idx: number) => (
                          <tr key={idx} className="border-b border-dashed">
                            <td className="py-1">{item.description}</td>
                            <td className="py-1 text-right">{item.amount.toLocaleString("ru-RU")} ₽</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                  <div className="text-right text-xs">
                    <p className="font-bold text-sm">Итого к оплате: {selectedBill.total.toLocaleString("ru-RU")} ₽</p>
                    <p className="text-[10px] text-muted-foreground mt-1">Оплатить в срок до: {new Date(selectedBill.due_date).toLocaleDateString("ru-RU")}</p>
                  </div>
                </div>
                <div className="flex justify-end gap-2 pt-2 border-t">
                  <Button variant="outline" onClick={() => { setShowInvoiceModal(false); setSelectedBill(null); }}>Закрыть</Button>
                  <Button onClick={() => window.print()} className="bg-slate-800 text-white font-bold">Напечатать счет</Button>
                </div>
              </>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
