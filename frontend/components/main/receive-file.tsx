import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";

import React from "react";

const ReceiveFilesCard = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Receive Files</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription>
          Receive files from other devices using the provided link.
        </CardDescription>
      </CardContent>
      <CardFooter>
        <Button>Receive Files</Button>
      </CardFooter>
    </Card>
  );
};

export default ReceiveFilesCard;
