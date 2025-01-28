import { Box, CssBaseline, AppBar, Toolbar, Typography, Drawer } from '@mui/material'
import { ReactNode } from 'react'
import Navbar from './Navbar'
import Sidebar from './Sidebar'

const drawerWidth = 240

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
    return (
      <Box sx={{ display: 'flex', minHeight: '100vh' }}>
        <CssBaseline />
        <Navbar />
        <Box
          component="nav"
          sx={{
            width: { sm: drawerWidth },
            flexShrink: { sm: 0 }
          }}
        >
          <Box
            sx={{
              position: 'fixed',
              width: drawerWidth,
              height: '100vh',
              bgcolor: 'background.paper',
              borderRight: (theme) => `1px solid ${theme.palette.divider}`,
            }}
          >
            <Toolbar /> {/* This creates space for the navbar */}
            <Sidebar />
          </Box>
        </Box>
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            p: 3,
            width: { sm: `calc(100% - ${drawerWidth}px)` },
            ml: { sm: `${drawerWidth}px` },
          }}
        >
          <Toolbar /> {/* This creates space for the navbar */}
          {children}
        </Box>
      </Box>
    );
  }