"use client";

import { useEffect, useState, useRef, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/use-user";
import { useProperties } from "@/hooks/use-properties";
import { ChatMessage } from "@/lib/types";
import { chatSocket } from "@/lib/chatSocket";
import { getProperties, getChatHistory } from "@/lib/properties";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export default function ChatPage() {
  const router = useRouter();
  const { data: user, isLoading } = useUser();

  const { data: rawProperties } = useProperties({ enabled: !!user });

  const properties = useMemo(() => {
    if (!user || !rawProperties) return [];
    if (user.role === "landlord") {
      return rawProperties.filter((p: any) => p.tenant);
    } else {
      return rawProperties.filter((p: any) => p.tenant?.tenantUserId === user.id || p.tenant?.tenant_user_id === user.id);
    }
  }, [user, rawProperties]);

  const [selectedProperty, setSelectedProperty] = useState<any | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputText, setInputText] = useState("");
  
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (user) {
      chatSocket.authenticate(user.id);
    }
  }, [user]);

  
  useEffect(() => {
    const unsubscribe = chatSocket.subscribe((incomingMsg) => {
      
      if (selectedProperty && incomingMsg.propertyId === selectedProperty.id) {
        setMessages((prev) => {
          
          if (prev.some(m => m.id === incomingMsg.id)) return prev;
          return [...prev, incomingMsg];
        });
      }
    });

    return () => unsubscribe();
  }, [selectedProperty]);

  
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/auth/login");
    }
  }, [isLoading, user, router]);

  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background text-foreground">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const handleSelectProperty = async (prop: any) => {
    setSelectedProperty(prop);
    try {
      const history = await getChatHistory(prop.id);
      const mappedHistory = history.map((m: any) => ({
        id: m.id,
        senderId: m.sender_id,
        receiverId: m.receiver_id,
        propertyId: m.property_id,
        text: m.text,
        timestamp: m.timestamp,
      }));
      setMessages(mappedHistory);
    } catch (e) {
      console.error("Failed to load message history", e);
    }
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim() || !selectedProperty || !selectedProperty.tenant) return;

    const tenant = selectedProperty.tenant as any;
    const prop = selectedProperty as any;
    const peerId = user.role === "landlord" 
      ? (tenant.tenantUserId || tenant.tenant_user_id)
      : (prop.landlordId || prop.landlord_id);

    const newMsg: ChatMessage = {
      id: "msg-" + Math.random().toString(36).substring(2, 9),
      senderId: user.id,
      receiverId: peerId,
      propertyId: selectedProperty.id,
      text: inputText.trim(),
      timestamp: new Date().toISOString(),
    };

    
    setMessages((prev) => [...prev, newMsg]);

    
    chatSocket.sendMessage(newMsg);

    setInputText("");
  };

  const isLandlord = user.role === "landlord";

  return (
    <div className="flex-1 pb-24 px-4 py-6 md:px-8 max-w-5xl mx-auto w-full flex h-[calc(100vh-80px)] max-h-[800px] gap-6">
      
      <div className="w-1/3 border border-slate-200 dark:border-slate-800 rounded-3xl bg-card shadow-xs flex flex-col overflow-hidden h-full">
        <div className="p-4 border-b">
          <h2 className="text-base font-black">Диалоги</h2>
          <p className="text-[10px] text-muted-foreground mt-0.5">Выберите собеседника для переписки</p>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-2">
          {properties.map((prop) => {
            const isSelected = selectedProperty?.id === prop.id;
            
            const pAny = prop as any;
            const peerName = isLandlord 
              ? prop.tenant?.name || "Арендатор"
              : pAny.landlordName || "Арендодатель";

            return (
              <button
                key={prop.id}
                onClick={() => handleSelectProperty(prop)}
                className={`w-full text-left p-3 rounded-2xl transition duration-150 flex flex-col ${
                  isSelected 
                    ? "bg-slate-100 dark:bg-slate-800" 
                    : "hover:bg-slate-50 dark:hover:bg-slate-900/50"
                }`}
              >
                <span className="font-bold text-xs truncate text-slate-800 dark:text-white">{peerName}</span>
                <span className="text-[10px] text-muted-foreground mt-0.5 truncate">{prop.name} ({prop.street})</span>
              </button>
            );
          })}

          {properties.length === 0 && (
            <div className="p-4 text-center text-xs text-muted-foreground italic">
              {isLandlord 
                ? "У вас нет сданных квартир с активными жильцами." 
                : "За вами не числится арендуемых апартаментов."}
            </div>
          )}
        </div>
      </div>

      
      <div className="flex-1 border border-slate-200 dark:border-slate-800 rounded-3xl bg-card shadow-xs flex flex-col overflow-hidden h-full">
        {selectedProperty && selectedProperty.tenant ? (
          <>
            
            <div className="p-4 border-b bg-slate-50/50 flex justify-between items-center">
              <div>
                <h3 className="text-xs font-bold text-slate-800 dark:text-white">
                  {isLandlord ? selectedProperty.tenant.name : (selectedProperty as any).landlordName || "Арендодатель"}
                </h3>
                <p className="text-[9px] text-muted-foreground mt-0.5">Квартира: {selectedProperty.name}</p>
              </div>
              <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" title="В сети (WebSocket)"></span>
            </div>

            
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {messages.map((msg) => {
                const isMyMessage = msg.senderId === user.id;
                
                return (
                  <div 
                    key={msg.id} 
                    className={`flex flex-col max-w-[70%] text-xs ${
                      isMyMessage ? "ml-auto items-end" : "mr-auto items-start"
                    }`}
                  >
                    <div className={`p-3 rounded-2xl shadow-xs leading-relaxed ${
                      isMyMessage 
                        ? "bg-slate-800 text-white rounded-br-xs" 
                        : "bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-white rounded-bl-xs border border-slate-200/50"
                    }`}>
                      <p>{msg.text}</p>
                    </div>
                    <span className="text-[9px] text-muted-foreground mt-1">
                      {new Date(msg.timestamp).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}
                    </span>
                  </div>
                );
              })}
              <div ref={messagesEndRef} />
            </div>

            
            <form onSubmit={handleSendMessage} className="p-4 border-t flex gap-2">
              <Input 
                placeholder="Напишите сообщение..." 
                className="h-10 text-xs rounded-xl flex-1"
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
              />
              <Button type="submit" className="h-10 rounded-xl bg-slate-800 hover:bg-slate-900 text-white font-semibold text-xs px-4">
                Отправить
              </Button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center text-center p-6 text-muted-foreground">
            <p className="text-sm font-semibold">Чат в реальном времени</p>
            <p className="text-xs max-w-xs mt-1">Выберите диалог из списка слева, чтобы начать общение с вашим жильцом или арендодателем.</p>
          </div>
        )}
      </div>
    </div>
  );
}
