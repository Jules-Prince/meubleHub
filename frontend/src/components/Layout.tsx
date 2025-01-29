import React from 'react';
import { Home, Box, Boxes, Menu, LogOut, User } from 'lucide-react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Sheet,
  SheetContent,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { authService } from '../services/auth';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const currentUser = authService.getCurrentUser();

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  const navigation = [
    { name: 'Main', href: '/', icon: Menu },
    { name: 'Home', href: '/home', icon: Home },
    { name: 'Room', href: '/room', icon: Box },
    { name: 'Object', href: '/object', icon: Boxes },
  ];

  const UserMenu = () => (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="flex items-center gap-2">
          <User className="h-5 w-5" />
          <div className="flex flex-col items-start">
            <span>{currentUser?.email}</span>
            {currentUser?.isAdmin && (
              <span className="text-xs text-red-500 font-semibold">Admin</span>
            )}
          </div>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem className="text-red-600" onClick={handleLogout}>
          <LogOut className="h-4 w-4 mr-2" />
          Logout
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );

  const Sidebar = () => (
    <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col lg:pt-16">
      <div className="flex flex-1 flex-col overflow-y-auto bg-gray-100 pt-5">
        <nav className="flex-1 space-y-1 px-2">
          {navigation.map((item) => (
            <Button
              key={item.name}
              variant={location.pathname === item.href ? "secondary" : "ghost"}
              className="w-full justify-start gap-2"
              onClick={() => navigate(item.href)}
            >
              <item.icon className="h-5 w-5" />
              {item.name}
            </Button>
          ))}
        </nav>
      </div>
    </div>
  );

  const MobileSidebar = () => (
    <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          className="lg:hidden"
        >
          <Menu className="h-6 w-6" />
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="w-64 p-0">
        <ScrollArea className="h-full px-2 py-4">
          <nav className="space-y-1">
            {navigation.map((item) => (
              <Button
                key={item.name}
                variant={location.pathname === item.href ? "secondary" : "ghost"}
                className="w-full justify-start gap-2"
                onClick={() => {
                  navigate(item.href);
                  setIsMobileMenuOpen(false);
                }}
              >
                <item.icon className="h-5 w-5" />
                {item.name}
              </Button>
            ))}
          </nav>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );

  return (
    <div className="min-h-screen">
      {/* Navbar */}
      <header className="fixed top-0 z-40 w-full border-b bg-background">
        <div className="flex h-16 items-center px-4">
          <MobileSidebar />
          <div className="flex flex-1 items-center justify-between">
            <h1 className="text-2xl font-bold">meubleHub</h1>
            {currentUser ? (
              <UserMenu />
            ) : (
              <Button onClick={() => navigate('/login')}>
                Login
              </Button>
            )}
          </div>
        </div>
      </header>

      {/* Sidebar */}
      <Sidebar />

      {/* Main Content */}
      <div className="lg:pl-64">
        <main className="pt-20 px-4">
          <div className="mx-auto max-w-7xl">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout;