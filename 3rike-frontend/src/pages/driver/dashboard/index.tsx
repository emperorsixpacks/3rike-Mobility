import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { useNavigate } from "react-router-dom";
import {
  ArrowUpRight,
  Plus,
  Hourglass,
  Home,
  Bell,
  User,
  Send,
  ArrowDown
} from "lucide-react";
import DepositModal from "../deposit";
import WithdrawOptions from "../withdraw/options";

export default function DriverDashboard() {
  const navigate = useNavigate();
  const [changeCurrency, setChangeCurrency] = useState(true);
  const [verificationStatus, setVerificationStatus] = useState<
    "not_started" | "in_progress" | "approved"
  >("not_started");

  // State to control the modals
  const [isDepositOpen, setIsDepositOpen] = useState(false);
  const [isWithdrawOpen, setIsWithdrawOpen] = useState(false);
  const [isMenuOpen, setIsMenuOpen] = useState(false); // State for the overlay menu

  useEffect(() => {
    const status = localStorage.getItem("verificationStatus") as
      | "not_started"
      | "in_progress"
      | "approved";

    if (status) {
      setVerificationStatus(status);
    }
  }, []);

  const handleLoan = () => {
    navigate('/driver/loan')
  };

  const handleSavings = () => {
    navigate('/driver/savings')
  };
  const handleInvestment = () => {
    navigate('/driver/investment')
  };
  const handleVerification = () => {
    navigate('/driver/verification')
  };
  const handle3rikeAi = () => {
    navigate('/driver/3rikeAi')
  };
  const handleNotification = () => {
    navigate('/driver/notification')
  };
  const handleSettings = () => {
    navigate('/driver/settings')
  };

  return (
    <div className="min-h-screen bg-white flex justify-center">
      {/* Mobile Frame Container */}
      <div className="w-full max-w-100 bg-white shadow-2xl overflow-hidden min-h-200 relative pb-10">

        {/* Header Profile */}
        <div className="px-6 flex items-center justify-between pt-6 mb-4">
          <div className="flex items-center gap-3 bg-white rounded-full">
            {/* Replace with actual user image */}
            <div className="w-10 h-10 rounded-full bg-gray-200 overflow-hidden border-2 border-white">
              <img src="/profile.png" alt="User" />
            </div>
            <span className="font-light text-sm -mr-5">Welcome, Effiong Musa</span>
            <Button variant="link">
              <img src="arrow.svg" alt="Arrow" className="w-5 h-5" />
            </Button>
          </div>
        </div>

        {/* Main Content Scroll Area */}
        <div className="px-5 space-y-4">

          {/* 1. GREEN BALANCE CARD */}
          <div className="relative w-full rounded-3xl p-6 text-white overflow-hidden">
            {/* Background Gradient & Blobs */}
            <img
              src="/earnings-banner.svg"
              alt="Card Background"
              className="absolute inset-0 w-full h-full bg-[#00C258] object-cover z-0"
            />

            <div className="relative z-10">
              <div className="flex justify-between items-start mb-2">
                <span className="text-green-100 text-sm font-light">Total Lifetime Earnings</span>
              </div>

              <div className="flex flex-row justify-between mb-6">
                <h1 className="text-4xl font-bold ">
                  {changeCurrency ? "$ 0.00" : "₵ 0.00"}
                </h1>

                {/* Custom Toggle Switch */}
                <div className="flex items-center rounded-full ">
                  <Switch
                    checked={changeCurrency}
                    onCheckedChange={setChangeCurrency}
                    className="data-[state=checked]:bg-black/25 data-[state=unchecked]:bg-black/25 h-6 w-10 border-5 border-transparent -rotate-90"
                  />
                </div>
              </div>


              <div className="flex gap-4">

                <Button
                  onClick={() => setIsDepositOpen(true)}
                  className="flex-1 bg-transparent hover:bg-white/30 text-white border border-white rounded-full h-12 gap-2 text-sm font-medium backdrop-blur-sm"
                >
                  <div className="bg-white text-[#00C258] rounded-full p-0.5 w-5 h-5 flex items-center justify-center">
                    <Plus size={14} strokeWidth={4} />
                  </div>
                  Deposit
                </Button>

                <Button onClick={() => setIsWithdrawOpen(true)} className="flex-1 bg-transparent hover:bg-white/30 text-white border border-white rounded-full h-12 gap-2 text-sm font-medium backdrop-blur-sm">
                  <div className="bg-white text-[#00C258] rounded-full p-0.5 w-5 h-5 flex items-center justify-center">
                    <ArrowUpRight size={14} strokeWidth={4} />
                  </div>
                  Withdraw
                </Button>
              </div>
            </div>
          </div>

          {/* 2. STATS ROW */}
          <div className="grid grid-cols-2 gap-4">
            {/* Savings Balance */}
            <div className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col gap-3">
              <div className="w-8 h-8 rounded-lg bg-purple-100 flex items-center justify-center">
                <img src="wallet.svg" alt="Back" className="w-6 h-6" />
              </div>
              <div>
                <p className="text-xs text-gray-400 mb-1">Savings Balance</p>
                <h3 className="text-xl font-light text-gray-800">$ 0.00</h3>
              </div>
            </div>

            {/* Pending Payout */}
            <div className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col gap-3 ">
              <div className="w-8 h-8 rounded-lg bg-orange-100 flex items-center justify-center">
                <Hourglass className="w-4 h-4 text-orange-500" />
              </div>
              <div>
                <p className="text-xs text-gray-400 mb-1">Pending Payout</p>
                <h3 className="text-xl font-light text-gray-800">$ 0.00</h3>
              </div>
            </div>
          </div>

          {/* to simulate approval */}
          <button
            type="button"
            onClick={() => setVerificationStatus("approved")}
            className="ml-2 text-xs text-[#1B8036] underline hover:[#1B8036] bg-transparent border-none cursor-pointer"
          >
            (Simulate Approval)
          </button>
          {/* 3. VERIFICATION BANNER */}
          <div onClick={handleVerification} className="relative w-full bg-[#1B8036] rounded-2xl p-5 overflow-hidden text-white flex items-center justify-between">

            {/* Abstract Background Pattern (Simulated with SVG) */}
            <img
              src="/verification-banner.svg"
              alt="Card Background"
              className="absolute inset-0 w-full h-full bg-[#1E8A32] object-cover z-0"
            />

            <div className="relative z-10">
              <div className="w-5 h-5 bg-white/20 rounded-lg flex items-center justify-center mb-2 backdrop-blur-md">
                {/* Icon simulating the scooter/delivery icon */}
                <img
                  src="/verification-bike.svg"
                  alt="Card Background"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </div>

              <h3 className="font-bold text-lg leading-tight">
                {verificationStatus === "not_started" && "Start Verification"}
                {verificationStatus === "in_progress" && "Verification in Progress"}
                {verificationStatus === "approved" && "Own a 3rike"}
              </h3>

              <p className="text-sm text-white mt-0.5">
                {verificationStatus === "not_started" && "Complete Kyc and be eligible."}
                {verificationStatus === "in_progress" && "Details received. We’ll be in touch soon."}
                {verificationStatus === "approved" && "Register and own a 3rike"}
              </p>
            </div>

            <div className="absolute bottom-4 right-4 z-10 w-8 h-8 bg-[#00C258] rounded-full flex items-center justify-center shadow-lg">
              <Button variant="link">
                <img
                  src="/arrow-right.svg"
                  alt="Card Background"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </Button>
            </div>
          </div>

          {/* 4. BOTTOM GRID MENU */}
          <div className="grid grid-cols-2 gap-4">
            {/* Savings */}
            <div onClick={handleSavings} className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col  gap-2 ">
              <div className="w-5 h-5 bg-white/20 rounded-lg flex items-center justify-center mb-2 backdrop-blur-md">
                <img
                  src="/piggy.svg"
                  alt="piggy"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </div>
              <span className="text-gray-800 text-lg">Savings</span>
            </div>

            {/* Investment */}
            <div onClick={handleInvestment} className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col gap-2 ">
              <div className="w-5 h-5 bg-white/20 rounded-lg flex items-center justify-center mb-2 backdrop-blur-md">
                <img
                  src="/invest.svg"
                  alt="invest"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </div>
              <span className="text-gray-800 text-lg">Investment</span>
            </div>

            {/* Earn */}
            <div className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col gap-2 ">
              <div className="w-5 h-5 bg-white/20 rounded-lg flex items-center justify-center mb-2 backdrop-blur-md">
                <img
                  src="/chart.svg"
                  alt="chart"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </div>
              <span className="text-gray-800 text-lg">Earn</span>
            </div>

            {/* Loan */}
            <div onClick={handleLoan} className="bg-white border-3 border-dashed border-gray-100 rounded-2xl p-4 flex flex-col gap-2 ">
              <div className="w-5 h-5 bg-white/20 rounded-lg flex items-center justify-center mb-2 backdrop-blur-md">
                <img
                  src="/loan.svg"
                  alt="loan"
                  className="absolute inset-0 w-full h-full  object-cover z-0"
                />
              </div>
              <span className="text-gray-800 text-lg">Loan</span>
            </div>
          </div>

        </div>

        {/* BOTTOM NAVIGATION BAR */}
        <div className="absolute w-full flex items-center justify-between px-6 -mt-7">
          {/* Left pill navigation */}
          <div className="bg-white rounded-full shadow-lg px-1 py-1 flex items-center -space-x-2">
            <Button
              variant="link"
              size="icon"
              className="hover:bg-transparent text-black"
            >
              <Home className="w-6 h-6 fill-current" />
            </Button>

            <Button
              onClick={handleNotification}
              variant="link"
              size="icon"
              className="hover:bg-transparent text-[#909090]"
            >
              <Bell className="w-6 h-6 fill-current" />
            </Button>

            <Button
              variant="link"
              onClick={handleSettings}
              size="icon"
              className="hover:bg-transparent text-gray-400"
            >
              <img
                src="/settings.svg"
                alt="settings"
                className=" w-5 h-5"
              />
            </Button>
          </div>

          {/* Right floating action button (TRIGGER) */}
          <Button
            variant="link"
            size="icon"
            onClick={() => setIsMenuOpen(true)} // Open modal on click
            className="hover:bg-transparent w-full h-full text-gray-400"
          >
            <img
              src="/add.svg"
              alt="add"
              className="ml-25 w-15 h-15"
            />
          </Button>
        </div>

        {/* --- OVERLAY MODAL --- */}
        {isMenuOpen && (
          <div className="absolute inset-0 z-50 bg-[#F3F5F9]/95 backdrop-blur-sm flex flex-col justify-end items-end p-2 animate-in fade-in duration-200">

            {/* Menu Items Container */}
            <div className="flex flex-col gap-6 mb-10 mr-20 items-start">

              {/* Option 1: Pay 3rike Ai */}
              <div onClick={handle3rikeAi} className="flex items-center gap-4 cursor-pointer group">
                <User className="w-6 h-6 text-[#00C259]" fill="#00C259" />
                <span className="text-lg font-light text-black group-hover:text-gray-700">Pay 3rike Ai</span>
              </div>

              {/* Option 2: Send */}
              <div className="flex items-center gap-4 cursor-pointer group">
                <Send className="w-6 h-6 text-[#9747FF]" fill="#9747FF" />
                <span className="text-lg font-light text-black group-hover:text-gray-700">Send</span>
              </div>

              {/* Option 3: Recieve */}
              <div className="flex items-center gap-4 cursor-pointer group">
                <ArrowDown className="w-6 h-6 text-[#FF9900]" strokeWidth={3} />
                <span className="text-lg font-light text-black group-hover:text-gray-700">Recieve</span>
              </div>
            </div>

            {/* Close Button (Replaces the Add Button position) */}
            <Button
              onClick={() => setIsMenuOpen(false)}
              className="w-25 h-25 pb-6 bg-transparent rounded-full flex items-end justify-end transition-transform hover:scale-105"
            >
              <img
                src="/subtract.svg"
                alt="subtract"
                className="w-15 h-15 -mr-2"
              />
            </Button>
          </div>
        )}

      </div>

      {/* deposit modal */}
      <DepositModal
        isOpen={isDepositOpen}
        onClose={() => setIsDepositOpen(false)}
      />

      {/* Withdraw modal */}
      <WithdrawOptions
        isOpen={isWithdrawOpen}
        onClose={() => setIsWithdrawOpen(false)}
      />
    </div>
  );
}