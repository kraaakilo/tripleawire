import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { useCallback, useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";
import { Button } from "@/components/ui/button.tsx";


export interface PacketMessage {
    source: string
    destination: string
    date: string
    protocol: string
    content: string
}


const Home = () => {
    const [socketUrl, setSocketUrl] = useState('wss://echo.websocket.org');
    const [packets, setPackets] = useState<PacketMessage[]>([]);
    const [alerts, setAlerts] = useState<string[]>([]);
    const [isCapturing, setIsCapturing] = useState(false);
    const { lastJsonMessage } = useWebSocket(socketUrl);

    const startCapturing = useCallback(() => {
        setIsCapturing(true);
        setSocketUrl('ws://localhost:8080/ws');
    }, [setIsCapturing, setSocketUrl]);

    const stopCapturing = useCallback(() => {
        setIsCapturing(false);
        setSocketUrl('wss://echo.websocket.org');
    }, [setIsCapturing, setSocketUrl]);

    useEffect(() => {
        if (lastJsonMessage && isCapturing) {

            let json = lastJsonMessage as { message: object | string, type: string };

            if (json.type === "packet") {
                const packet: PacketMessage = json.message as PacketMessage;
                setPackets(prevPackets => [packet, ...prevPackets]);
            } else if (json.type === "alert") {
                const alert = json.message as string;
                if (!alerts.includes(alert)) {
                    setAlerts(prevAlerts => [alert, ...prevAlerts]);
                }
            }
        }
    }, [lastJsonMessage]);

    return (
        <div className="container my-12">
            <div className="flex items-center justify-between">
                <h1 className="flex gap-2 text-4xl font-bold items-center">
                    <img src="/icon.svg" alt="Icons" className="size-12" />
                    TripleaWire Packets Viewer
                </h1>
                {!isCapturing && <Button onClick={startCapturing}>Start Capturing !</Button>}
                {isCapturing && <Button onClick={stopCapturing} variant="destructive">Stop Capturing !</Button>}

            </div>
            <div className="grid gap-5 mt-6">
                <div>
                    <h2 className="text-2xl font-bold">Packets</h2>
                    <div className="h-96 overflow-y-scroll ">
                        <Table className="mt-6">
                            <TableCaption>
                                All received packets are listed below.
                            </TableCaption>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[100px]">Protocol</TableHead>
                                    <TableHead>Source IP</TableHead>
                                    <TableHead>Destination IP</TableHead>
                                    <TableHead className="text-right">Payload</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {packets.map((packet, index) => (
                                    <TableRow key={index}>
                                        <TableCell>{packet.protocol}</TableCell>
                                        <TableCell>{packet.source}</TableCell>
                                        <TableCell>{packet.destination}</TableCell>
                                        <TableCell className="text-right">{packet.content}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </div>
                </div>
                <div className="h-96">
                    <h2 className="text-2xl font-bold">Alerts</h2>
                    <ul className="mt-4">
                        {alerts.map((alert, index) => (
                            <li key={index} className="p-2 bg-red-100 text-red-900 rounded-md">{alert}</li>
                        ))}
                    </ul>
                </div>
            </div>

        </div>
    );
};

export default Home;