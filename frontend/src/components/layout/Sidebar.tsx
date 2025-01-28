import {
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Divider,
  } from '@mui/material';
  import { useNavigate, useLocation } from 'react-router-dom';
  import HomeIcon from '@mui/icons-material/Home';
  import MeetingRoomIcon from '@mui/icons-material/MeetingRoom';
  import ChairIcon from '@mui/icons-material/Chair';
  import DashboardIcon from '@mui/icons-material/Dashboard';
  
  interface NavigationItem {
    text: string;
    path: string;
    icon: JSX.Element;
  }
  
  const navigationItems: NavigationItem[] = [
    {
      text: 'Dashboard',
      path: '/',
      icon: <DashboardIcon />,
    },
    {
      text: 'Homes',
      path: '/homes',
      icon: <HomeIcon />,
    },
    {
      text: 'Rooms',
      path: '/rooms',
      icon: <MeetingRoomIcon />,
    },
    {
      text: 'Objects',
      path: '/objects',
      icon: <ChairIcon />,
    },
  ];
  
  export default function Sidebar() {
    const navigate = useNavigate();
    const location = useLocation();
  
    return (
      <>
        <Divider />
        <List>
          {navigationItems.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => navigate(item.path)}
              >
                <ListItemIcon sx={{
                  color: location.pathname === item.path ? 'primary.main' : 'inherit'
                }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText 
                  primary={item.text}
                  sx={{
                    color: location.pathname === item.path ? 'primary.main' : 'inherit'
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </>
    );
  }