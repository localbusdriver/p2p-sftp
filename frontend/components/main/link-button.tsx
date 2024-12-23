"use client";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";


const LinkButton = ({
  children,
  link,
  className,
  variant,
}: {
  children: React.ReactNode;
  link: string;
  className?: string;
  variant?:
    | "link"
    | "default"
    | "destructive"
    | "outline"
    | "secondary"
    | "ghost"
    | null
    | undefined;
}) => {
  const router = useRouter();

  return (
    <Button
      className={className}
      onClick={() => router.push(link)}
      variant={variant}
    >
      {children}
    </Button>
  );
};

export default LinkButton;
