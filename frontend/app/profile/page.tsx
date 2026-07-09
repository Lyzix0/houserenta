"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/use-user";
import { logout, updateProfile } from "@/lib/auth";
import { useQueryClient } from "@tanstack/react-query";
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
  DialogTrigger,
} from "@/components/ui/dialog";
import { AuthGuard } from "@/components/auth-guard";
import { HugeiconsIcon } from "@hugeicons/react";
import { Edit03Icon, LogoutSquare01Icon, UserIcon } from "@hugeicons/core-free-icons";

export default function ProfilePage() {
  return (
    <AuthGuard mode="auth">
      <ProfileContent />
    </AuthGuard>
  );
}

function ProfileContent() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: user, isLoading, refetch } = useUser();

  const [name, setName] = useState("");
  const [document, setDocument] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [paymentCard, setPaymentCard] = useState("");
  const [password, setPassword] = useState("");

  const [editError, setEditError] = useState("");
  const [editSuccess, setEditSuccess] = useState("");
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    if (user) {
      setName(user.name || "");
      setDocument(user.document || "");
      setPhone(user.phone || "");
      setEmail(user.email || "");
      setPaymentCard(user.paymentCard || "");
    }
  }, [user]);

  useEffect(() => {
    if (isClient && !isLoading && !user) {
      router.push("/auth/login");
    }
  }, [isClient, isLoading, user, router]);

  if (isLoading || !isClient) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background text-foreground">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setEditError("");
    setEditSuccess("");

    try {
      await updateProfile({
        name,
        document,
        phone,
        paymentCard: user.role === "landlord" ? paymentCard : undefined,
        email,
        password: password ? password : undefined,
      });
      setEditSuccess("Профиль успешно обновлен!");
      setPassword("");
      refetch();
      setTimeout(() => {
        setIsEditDialogOpen(false);
        setEditSuccess("");
      }, 1500);
    } catch (err: any) {
      setEditError(err.message || "Ошибка сохранения");
    }
  };

  const handleLogout = async () => {
    await logout();
    queryClient.clear();
    router.push("/auth/login");
  };

  return (
    <div className="flex-1 pb-24 px-4 py-6 md:px-8 max-w-xl mx-auto w-full">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">Профиль</h1>
          </div>
          <p className="text-xs text-muted-foreground mt-0.5">Личные настройки и выбор активной роли</p>
        </div>
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogTrigger asChild>
            <Button 
              variant="ghost" 
              size="icon-lg"
              className="size-11 bg-muted hover:bg-neutral-400/40"
            >
              <HugeiconsIcon icon={Edit03Icon} className="size-5" />
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Редактирование профиля</DialogTitle>
              <DialogDescription>
                Внесите изменения в ваши личные данные. Оставьте поле пароля пустым, если не хотите его менять.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleUpdateProfile} className="space-y-4 pt-2">
              <div className="space-y-2">
                <Label htmlFor="profName">ФИО *</Label>
                <Input 
                  id="profName" 
                  value={name} 
                  onChange={e => setName(e.target.value)} 
                  required 
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="profDoc">Документ (паспорт)</Label>
                <Input 
                  id="profDoc" 
                  value={document} 
                  onChange={e => setDocument(e.target.value)} 
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="profPhone">Телефон *</Label>
                <Input 
                  id="profPhone" 
                  value={phone} 
                  onChange={e => setPhone(e.target.value)} 
                  required 
                />
              </div>

              {user.role === "landlord" && (
                <div className="space-y-2">
                  <Label htmlFor="profCard">Карта для получения оплаты</Label>
                  <Input 
                    id="profCard" 
                    placeholder="Карта не указана"
                    value={paymentCard} 
                    onChange={e => setPaymentCard(e.target.value)} 
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="profEmail">Почта / Логин *</Label>
                <Input 
                  id="profEmail" 
                  type="email"
                  value={email} 
                  onChange={e => setEmail(e.target.value)} 
                  required 
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="profPass">Новый пароль</Label>
                <Input 
                  id="profPass" 
                  type="password"
                  placeholder="Введите новый пароль (не менее 6 симв.)"
                  value={password} 
                  onChange={e => setPassword(e.target.value)} 
                />
              </div>

              {editSuccess && (
                <p className="text-xs text-green-600 font-semibold text-center">{editSuccess}</p>
              )}
              {editError && (
                <p className="text-xs text-red-500 font-semibold text-center">{editError}</p>
              )}

              <div className="flex gap-2 justify-end pt-2">
                <Button type="button" variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                  Отмена
                </Button>
                <Button type="submit">
                  Сохранить изменения
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
        <Button 
          onClick={handleLogout} 
          variant="destructive" 
          size="icon-lg"
          className="size-11 p-0"
        >
          <HugeiconsIcon icon={LogoutSquare01Icon} className="size-5"/>
        </Button>
      </div>
      <div className="px-2 py-2 mb-4 flex justify-start items-center gap-2">
        <div className="bg-muted rounded-full p-2 border">
          <HugeiconsIcon icon={UserIcon} size={45}/>
        </div>
        <div className="px-2">
          <h2 className="text-[22px] font-medium">{name}</h2>
          <div className="flex gap-1.5 items-center mt-0.5">
            <span className="text-sm text-foreground/60">{email}</span>
            <span className="text-xs text-foreground/60">•</span>
            <Badge variant="secondary" className="text-xs py-0.5 font-semibold">
              {user.role === "landlord" ? "Арендодатель" : "Арендатор"}
            </Badge>
          </div>
        </div>
      </div>
      <Card className="border-border shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-bold">Личные данные</CardTitle>
          <CardDescription>Контактная информация и реквизиты вашего аккаунта</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-xs md:text-sm">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-b border-dashed border-border/50 pb-3">
            <div>
              <p className="text-muted-foreground font-medium">Полное имя</p>
              <p className="font-semibold text-foreground mt-0.5">{user.name || "Не указано"}</p>
            </div>
            <div>
              <p className="text-muted-foreground font-medium">Документ (паспорт)</p>
              <p className="font-semibold text-foreground mt-0.5">{user.document || "Не указан"}</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 border-b border-dashed border-border/50 pb-3">
            <div>
              <p className="text-muted-foreground font-medium">Телефон</p>
              <p className="font-semibold text-foreground mt-0.5">{user.phone || "Не указан"}</p>
            </div>
            <div>
              <p className="text-muted-foreground font-medium">Почта / Логин</p>
              <p className="font-semibold text-foreground mt-0.5">{user.email || "Не указана"}</p>
            </div>
          </div>

          {user.role === "landlord" && (
            <div className="pb-1">
              <p className="text-muted-foreground font-medium">Карта для получения оплаты</p>
              <p className="font-semibold text-foreground mt-0.5">{user.paymentCard || "Карта не указана"}</p>
              <p className="text-[10px] text-muted-foreground mt-1">Отображается в печатных квитанциях для арендатора.</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
