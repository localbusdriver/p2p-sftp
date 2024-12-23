"use client";

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { GetUsername, SetUsername } from "@/wailsjs/go/app/App";
import { useRouter } from "next/navigation";

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const formSchema = z.object({
  username: z
    .string()
    .min(3, {
      message: "Username must be at least 3 characters long",
    })
    .max(20, {
      message: "Username must be at most 20 characters long",
    }),
});

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
    },
  });

  useEffect(() => {
    const loadUsername = async () => {
      try {
        const username = await GetUsername();
        if (username) {
          form.setValue("username", username);
        }
      } catch (error) {
        console.error("Failed to load username", error);
      }
    };
    setIsLoading(true);
    loadUsername();
    setIsLoading(false);
  }, [form]);

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsLoading(true);
    setError(null);
    setSuccess(false);

    try {
      await SetUsername(values.username);
      setSuccess(true);
      form.reset();
      router.push("/main");
      console.log("Username saved successfully!");
    } catch (error) {
      console.error("Failed to save username", error);
      setError(
        error instanceof Error ? error.message : "Failed to save username"
      );
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <div>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="mt-8 sm:mt-16 md:mt-24 lg:mt-36 xl:mt-52 w-[300px] flex items-start justify-center flex-col gap-2 mx-auto"
        >
          <h1 className="font-bold">Username</h1>
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem className="flex flex-col gap-2 items-start justify-center">
                <div className="flex flex-row gap-2 items-center justify-center">
                  <FormControl>
                    <Input
                      placeholder="12345"
                      {...field}
                      disabled={isLoading}
                      className=""
                    />
                  </FormControl>
                  {isLoading ? (
                    <span className="text-gray-500">Loading... </span>
                  ) : (
                    <Button type="submit">Submit</Button>
                  )}
                </div>
                <FormDescription>
                  {error ? (
                    <span className="text-red-500">{error}</span>
                  ) : success ? (
                    <span className="text-green-500">
                      Username saved successfully!
                    </span>
                  ) : (
                    "Username for your peers to find you!"
                  )}
                </FormDescription>
              </FormItem>
            )}
          />
          <FormLabel>Continue with your username</FormLabel>
        </form>
      </Form>
    </div>
  );
}
