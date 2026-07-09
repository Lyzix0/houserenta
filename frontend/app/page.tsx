import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  Building2,
  ReceiptText,
  BellRing,
  Users,
  Gauge,
  Wallet,
  CheckCircle2,
  ArrowRight,
  AlertTriangle,
  CalendarClock,
  FileText,
  Smartphone,
} from "lucide-react";

function DashboardMockup() {
  const chartBars = [42, 58, 50, 72, 64, 88, 80, 96];

  return (
    <div className="relative">
      <div
        aria-hidden
        className="absolute -inset-6 rounded-3xl bg-primary/20 blur-3xl"
      />
      <div className="relative overflow-hidden rounded-2xl border bg-card text-card-foreground shadow-2xl">
        <div className="flex items-center gap-2 border-b bg-muted/50 px-4 py-3">
          <span className="h-2.5 w-2.5 rounded-full bg-muted-foreground/30" />
          <span className="h-2.5 w-2.5 rounded-full bg-muted-foreground/30" />
          <span className="h-2.5 w-2.5 rounded-full bg-muted-foreground/30" />
          <span className="ml-3 text-xs font-medium text-muted-foreground">
            Моя Аренда — главная
          </span>
        </div>

        <div className="grid gap-4 p-4 sm:grid-cols-5">
          <div className="space-y-3 sm:col-span-3">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Требует внимания
            </p>

            <div className="flex items-start gap-3 rounded-xl border border-destructive/20 bg-destructive/10 p-3">
              <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-destructive" />
              <div>
                <p className="text-sm font-medium">Отрицательный баланс</p>
                <p className="text-xs text-muted-foreground">
                  Ленина 12, кв. 45 · −4 320 ₽
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3 rounded-xl border bg-muted/60 p-3">
              <CalendarClock className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Скорый выезд жильца</p>
                <p className="text-xs text-muted-foreground">
                  Садовая 8, кв. 12 · через 6 дней
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3 rounded-xl border border-primary/20 bg-primary/10 p-3">
              <Gauge className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
              <div>
                <p className="text-sm font-medium">Дедлайн по показаниям</p>
                <p className="text-xs text-muted-foreground">
                  Мира 3, кв. 7 · осталось 2 дня
                </p>
              </div>
            </div>
          </div>

          <div className="space-y-3 sm:col-span-2">
            <div className="rounded-xl border p-3">
              <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Доход за 8 мес.
              </p>
              <div className="mt-3 flex h-24 items-end gap-1.5">
                {chartBars.map((h, i) => (
                  <div
                    key={i}
                    style={{ height: `${h}%` }}
                    className={`flex-1 rounded-t ${
                      i === chartBars.length - 1
                        ? "bg-chart-1"
                        : "bg-chart-1/30"
                    }`}
                  />
                ))}
              </div>
            </div>

            <div className="rounded-xl border p-3">
              <div className="flex items-center justify-between gap-2">
                <p className="text-sm font-medium">Счёт за октябрь</p>
                <Badge>выставлен</Badge>
              </div>
              <p className="mt-1 text-xs text-muted-foreground">
                Сформирован автоматически
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function TenantPhoneMockup() {
  return (
    <div className="relative mx-auto w-64">
      <div
        aria-hidden
        className="absolute -inset-8 rounded-full bg-primary/15 blur-3xl"
      />
      <div className="relative rounded-[2rem] border-8 border-foreground bg-card text-card-foreground shadow-2xl">
        <div className="mx-auto mt-2 h-1.5 w-16 rounded-full bg-muted" />
        <div className="space-y-3 p-4 pb-6">
          <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Кабинет жильца
          </p>
          <div className="rounded-xl bg-primary p-4 text-primary-foreground">
            <p className="text-xs opacity-80">Баланс</p>
            <p className="text-2xl font-bold">+2 150 ₽</p>
          </div>
          <div className="flex items-center justify-between rounded-xl border p-3">
            <div className="flex items-center gap-2">
              <Gauge className="h-4 w-4 text-primary" />
              <span className="text-sm">Передать показания</span>
            </div>
            <ArrowRight className="h-4 w-4 text-muted-foreground/50" />
          </div>
          <div className="flex items-center justify-between rounded-xl border p-3">
            <div className="flex items-center gap-2">
              <FileText className="h-4 w-4 text-primary" />
              <span className="text-sm">Счета и история</span>
            </div>
            <ArrowRight className="h-4 w-4 text-muted-foreground/50" />
          </div>
          <Button className="w-full">Оплатить счёт</Button>
        </div>
      </div>
    </div>
  );
}

const features = [
  {
    icon: ReceiptText,
    title: "Автоматические счета",
    text: "Система сама формирует квитанции каждые 30 дней на основе переданных показаний и ваших тарифов.",
  },
  {
    icon: BellRing,
    title: "Система предупреждений",
    text: "Главная страница выводит только важное: уведомления об отрицательном балансе, скором выезде жильца и дедлайнах по показаниям.",
  },
  {
    icon: Users,
    title: "Учёт объектов и жильцов",
    text: "Добавляйте квартиры, настраивайте тарифы ЖКХ и заводите карточки арендаторов в пару кликов. Вся история начислений и оплат всегда под рукой.",
  },
];

const pains = [
  "Ошибки в цифрах при ручных расчётах",
  "Забытые показания счётчиков",
  "Личные напоминания жильцам об оплате",
];

const steps = [
  {
    title: "Зарегистрируйтесь",
    text: "Создайте аккаунт и добавьте свои объекты недвижимости.",
  },
  {
    title: "Настройте тарифы",
    text: "Укажите тарифы ЖКХ и заведите карточку арендатора.",
  },
  {
    title: "Передайте доступ жильцу",
    text: "Отправьте логин и пароль — дальше система работает автоматически.",
  },
];

const tenantBenefits = [
  { icon: Wallet, text: "Баланс и история оплат всегда на виду" },
  { icon: Gauge, text: "Передача показаний за 30 секунд" },
  { icon: ReceiptText, text: "Оплата счёта в один клик" },
];

export default function Home() {
  return (
    <main className="min-h-screen bg-background text-foreground antialiased">
      <header className="sticky top-0 z-50 border-b bg-background/80 backdrop-blur-md">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
          <Link href="/" className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
              <Building2 className="h-5 w-5 text-primary-foreground" />
            </span>
            <span className="text-lg font-bold tracking-tight">Моя Аренда</span>
          </Link>
          <nav className="hidden items-center gap-8 text-sm font-medium text-muted-foreground md:flex">
            <a href="#features" className="transition-colors hover:text-primary">
              Возможности
            </a>
            <a href="#tenant" className="transition-colors hover:text-primary">
              Для жильцов
            </a>
            <a href="#how" className="transition-colors hover:text-primary">
              Как начать
            </a>
          </nav>
          <Button asChild>
            <Link href="/auth/signup">Регистрация</Link>
          </Button>
        </div>
      </header>

      <section className="relative overflow-hidden">
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0 bg-[radial-gradient(60rem_30rem_at_70%_-10%,hsl(var(--primary)/0.12),transparent)]"
        />
        <div className="relative mx-auto grid max-w-6xl items-center gap-14 px-4 pb-20 pt-16 lg:grid-cols-2 lg:pb-28 lg:pt-24">
          <div>
            <Badge variant="outline" className="mb-5 border-primary/30 bg-primary/10 text-primary">
              Профессиональный инструмент арендодателя
            </Badge>
            <h1 className="text-4xl font-extrabold leading-[1.1] tracking-tight sm:text-5xl">
              Полный контроль над сдачей недвижимости{" "}
              <span className="text-primary">в одном приложении</span>
            </h1>
            <p className="mt-6 max-w-lg text-lg leading-relaxed text-muted-foreground">
              Автоматизируйте расчёты, выставляйте счета вовремя и забудьте о
              рутине. Инструмент, который экономит время и исключает финансовые
              ошибки.
            </p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row">
              <Button size="lg" asChild className="text-base">
                <Link href="/auth/signup">
                  Зарегистрироваться
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild className="text-base">
                <a href="#how">Как начать</a>
              </Button>
            </div>
            <div className="mt-8 flex flex-wrap gap-x-6 gap-y-2 text-sm text-muted-foreground">
              <span className="flex items-center gap-1.5">
                <CheckCircle2 className="h-4 w-4 text-primary" />
                Счета каждые 30 дней автоматически
              </span>
              <span className="flex items-center gap-1.5">
                <CheckCircle2 className="h-4 w-4 text-primary" />
                Без ошибок в расчётах
              </span>
            </div>
          </div>
          <DashboardMockup />
        </div>
      </section>

      <section className="border-y bg-muted/40">
        <div className="mx-auto grid max-w-6xl items-center gap-10 px-4 py-16 lg:grid-cols-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">
              Забудьте о ручных расчётах
            </h2>
            <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
              Больше никаких ошибок в цифрах, забытых показаний и необходимости
              лично напоминать жильцам об оплате.{" "}
              <span className="font-semibold text-foreground">Моя Аренда</span>{" "}
              берёт всю рутину на себя.
            </p>
          </div>
          <ul className="space-y-3">
            {pains.map((pain) => (
              <li
                key={pain}
                className="flex items-center gap-3 rounded-xl border bg-card p-4 text-card-foreground"
              >
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                  <CheckCircle2 className="h-4 w-4 text-primary" />
                </span>
                <span>
                  {pain} —{" "}
                  <span className="font-medium text-primary">решено</span>
                </span>
              </li>
            ))}
          </ul>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-6xl px-4 py-20 lg:py-24">
        <p className="text-center text-sm font-semibold uppercase tracking-widest text-primary">
          Возможности
        </p>
        <h2 className="mt-3 text-center text-3xl font-bold tracking-tight sm:text-4xl">
          Ключевые возможности для арендодателя
        </h2>
        <div className="mt-14 grid gap-6 md:grid-cols-3">
          {features.map(({ icon: Icon, title, text }) => (
            <Card
              key={title}
              className="group transition-all hover:-translate-y-1 hover:border-primary/30 hover:shadow-lg hover:shadow-primary/10"
            >
              <CardHeader>
                <div className="mb-3 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 transition-colors group-hover:bg-primary">
                  <Icon className="h-6 w-6 text-primary transition-colors group-hover:text-primary-foreground" />
                </div>
                <CardTitle className="text-lg">{title}</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-sm leading-relaxed">
                  {text}
                </CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section id="tenant" className="border-y bg-muted/40">
        <div className="mx-auto grid max-w-6xl items-center gap-14 px-4 py-20 lg:grid-cols-2 lg:py-24">
          <TenantPhoneMockup />
          <div className="order-first lg:order-none">
            <p className="flex items-center gap-2 text-sm font-semibold uppercase tracking-widest text-primary">
              <Smartphone className="h-4 w-4" />
              Удобство для арендатора
            </p>
            <h2 className="mt-4 text-3xl font-bold tracking-tight sm:text-4xl">
              Жильцы получают понятный личный кабинет
            </h2>
            <p className="mt-5 text-lg leading-relaxed text-muted-foreground">
              Они могут просматривать баланс, передавать показания и оплачивать
              счета в один клик. Это гарантирует прозрачность и исключает
              недопонимание между собственником и жильцом.
            </p>
            <ul className="mt-7 space-y-3">
              {tenantBenefits.map(({ icon: Icon, text }) => (
                <li key={text} className="flex items-center gap-3 text-muted-foreground">
                  <Icon className="h-5 w-5 shrink-0 text-primary" />
                  {text}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </section>

      <section id="how" className="mx-auto max-w-6xl px-4 py-20 lg:py-24">
        <h2 className="text-center text-3xl font-bold tracking-tight sm:text-4xl">
          Как начать
        </h2>
        <div className="relative mt-14 grid gap-10 md:grid-cols-3">
          <div
            aria-hidden
            className="absolute left-0 right-0 top-6 hidden h-px bg-gradient-to-r from-transparent via-border to-transparent md:block"
          />
          {steps.map(({ title, text }, i) => (
            <div key={title} className="relative text-center md:text-left">
              <div className="relative z-10 mx-auto mb-5 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-lg font-bold text-primary-foreground shadow-lg shadow-primary/30 md:mx-0">
                {i + 1}
              </div>
              <h3 className="mb-2 text-lg font-semibold">{title}</h3>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {text}
              </p>
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-4 pb-24">
        <div className="relative overflow-hidden rounded-3xl bg-primary px-6 py-16 text-center text-primary-foreground md:px-16">
          <div
            aria-hidden
            className="pointer-events-none absolute inset-0 bg-[radial-gradient(40rem_20rem_at_50%_120%,hsl(var(--primary-foreground)/0.15),transparent)]"
          />
          <h2 className="relative text-3xl font-bold tracking-tight sm:text-4xl">
            Управляйте недвижимостью профессионально
          </h2>
          <p className="relative mx-auto mt-4 max-w-xl text-lg opacity-90">
            Освободите своё время от рутины и возьмите контроль над финансами в
            свои руки.
          </p>
          <Button
            size="lg"
            variant="secondary"
            asChild
            className="relative mt-9 text-base font-semibold"
          >
            <Link href="/auth/signup">
              Зарегистрироваться
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </section>

      <footer className="border-t">
        <div className="mx-auto max-w-6xl px-4 py-8">
          <div className="flex flex-col items-center justify-between gap-4 text-sm text-muted-foreground md:flex-row">
            <div className="flex items-center gap-2">
              <span className="flex h-6 w-6 items-center justify-center rounded-md bg-primary">
                <Building2 className="h-3.5 w-3.5 text-primary-foreground" />
              </span>
              <span>Моя Аренда © {new Date().getFullYear()}</span>
            </div>
            <Separator className="md:hidden" />
            <div className="flex gap-6">
              <a href="#features" className="transition-colors hover:text-primary">
                Возможности
              </a>
              <a href="#how" className="transition-colors hover:text-primary">
                Как начать
              </a>
              <Link href="/auth/signup" className="transition-colors hover:text-primary">
                Регистрация
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </main>
  );
}