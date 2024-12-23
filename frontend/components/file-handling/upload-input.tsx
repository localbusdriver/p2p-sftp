import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { UploadFile } from "@/wailsjs/go/app/App";
import { useState } from "react";
import { makeid } from "@/lib/utils";

const UploadInput = () => {
  const [fileName, setFileName] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const handleFileUpload = async () => {
    if (!file) {
      return;
    }

    if (!fileName) {
      setFileName(makeid(10));
    }

    try {
      const arrayBuffer = await file.arrayBuffer();
      const unit8Array = new Uint8Array(arrayBuffer);
      const normArray = Array.from(unit8Array);

      await UploadFile(fileName, normArray);
    } catch (e) {
      console.error("Failed to upload file", e);
    }
  };

  const handleChangeFileState = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) {
      return;
    }

    setFile(e.target.files[0]);
    console.log(fileName)
    if (!fileName) {
      const extension = e.target.files[0].name.match(/\.(pdf|jpg|png)$/)?.[0];
      const noExtFileName = e.target.files[0].name.replace(/\.(pdf|jpg|png)$/, "");
      setFileName(`${noExtFileName}_${makeid(10)}${extension}`);
    }
  };

  const handleClear = () => {
    setFileName("");
    console.log(fileName)
    setFile(null);
  };

  return (
    <div className="flex flex-col gap-3 w-full sm:w-[360px] md:w-[400px] px-4 py-3 border rounded-md shadow-md"> 
      <div className="w-full grid grid-cols-11 items-center">
        <Label htmlFor="file-name" className="col-span-5">
          File name:
        </Label>
        <Input
          id="file-name"
          type="text"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setFileName(e.target.value)
          }
          value={fileName}
          className="w-full col-span-6"
        />
      </div>
      <div className="w-full grid grid-cols-11 items-center">
        <Label htmlFor="file-upload" className="col-span-5">
          Upload File:
        </Label>

        <Input
          id="file-upload"
          type="file"
          accept="image/png, image/jpeg, image/jpg, application/pdf"
          onChange={handleChangeFileState}
          className="w-full col-span-6"
        />
      </div>
      <div className="flex gap-2 items-center ">
        <Button variant="outline" onClick={handleClear} className="w-full">
          Clear
        </Button>
        <Button onClick={handleFileUpload} className="w-full">
          Store
        </Button>
      </div>
    </div>
  );
};

export default UploadInput;
