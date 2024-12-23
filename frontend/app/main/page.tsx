"use client";
import { GetUsername, GetUserId } from "@/wailsjs/go/app/App";
import {Button} from "@/components/ui/button";
import { useState, useEffect } from "react";
import { useToast } from "@/hooks/use-toast";
import UploadInput from "@/components/file-handling/upload-input";

const Page = () => {
  const {toast} = useToast();
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    const loadUsername = async () => {
      try {
        const configUsername = await GetUsername();
        setUsername(configUsername);
      } catch (error) {
        console.error("Failed to load username", error);
      }
    };
    loadUsername();
  }, []);

  const loadUserId = async () => {
    const fetchedUserId = await GetUserId();
    if (!fetchedUserId) {
      return;
    }
    const copiedToast = (iteration: number = 0) => {
      const title = `Copied to clipboard! x${iteration + 1}`;
      toast({
        title: title,
        description: (
          <Button
            className="bg-amber-500 text-white"
            onClick={() => {
              navigator.clipboard.writeText(fetchedUserId);
              copiedToast(iteration + 1);
            }}
          >
            Copy again&nbsp;
            {iteration > 9 ? "~ ∞" : iteration > 0 ? `x${iteration + 1}`  : ""}
          </Button>
        ),
        variant: "default",
      });
    }
    toast({
      title: fetchedUserId,
      description: <Button className="bg-blue-500 text-white" onClick={() => {
        navigator.clipboard.writeText(fetchedUserId)
        copiedToast()
      }}>User ID</Button>,
      variant: "default",
    });
  }

  return (
    <div className="flex flex-col items-center justify-center gap-6">
      <h1 className="text-center">
        Hello, <span className="font-bold">{username && username}!</span>
      </h1>

      <Button onClick={loadUserId}>Show ID</Button>
      <UploadInput />
    </div>
  );
};

export default Page;
