import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";

const SendFileCard = () => {
  return <Card>
    <CardHeader>
      <CardTitle>Send a file</CardTitle>
    </CardHeader>
    <CardContent> 
      <CardDescription>
        Open a SFTP server and send a file to your peers!
      </CardDescription>
    </CardContent>
    <CardFooter>
      <Button variant="outline">Send a file</Button>
      </CardFooter>
  </Card>;
};

export default SendFileCard;
